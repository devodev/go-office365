package office365

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Watcher is an interface used by Watch for generating a stream of records.
type Watcher interface {
	Run(context.Context) chan Resource
}

// SubscriptionWatcher implements the Watcher interface.
// It fecthes current subscriptions, then queries content available for a given interval
// and proceed to query audit records.
type SubscriptionWatcher struct {
	client *Client
	config SubscriptionWatcherConfig

	// message bus
	queue chan Resource

	// state
	muBusy          *sync.Mutex
	contentTypeBusy map[ContentType]bool

	State
}

// SubscriptionWatcherConfig .
type SubscriptionWatcherConfig struct {
	LookBehindMinutes     int
	TickerIntervalSeconds int
}

// NewSubscriptionWatcher returns a new watcher that uses the provided client
// for querying the API.
func NewSubscriptionWatcher(client *Client, conf SubscriptionWatcherConfig, s State) (*SubscriptionWatcher, error) {
	lookBehindDur := time.Duration(conf.LookBehindMinutes) * time.Minute
	if lookBehindDur <= 0 {
		return nil, fmt.Errorf("lookBehindMinutes must be greater than 0")
	}
	if lookBehindDur > 24*time.Hour {
		return nil, fmt.Errorf("lookBehindMinutes must be less than or equal to 24 hours")
	}

	tickerIntervalDur := time.Duration(conf.TickerIntervalSeconds) * time.Second
	if tickerIntervalDur <= 0 {
		return nil, fmt.Errorf("tickerIntervalSeconds must be greater than 0")
	}
	if tickerIntervalDur > time.Hour {
		return nil, fmt.Errorf("tickerIntervalSeconds must be less than or equal to 1 hour")
	}

	watcher := &SubscriptionWatcher{
		client: client,
		config: conf,

		queue: make(chan Resource, contentTypeCount),

		muBusy:          &sync.Mutex{},
		contentTypeBusy: make(map[ContentType]bool),
		State:           s,
	}
	return watcher, nil
}

func (s *SubscriptionWatcher) sendResourceOrSkip(ctx context.Context, r Resource) {
	select {
	default:
		return
	case <-ctx.Done():
		return
	case s.queue <- r:
	}
}

func (s *SubscriptionWatcher) isBusy(ct *ContentType) bool {
	s.muBusy.Lock()
	defer s.muBusy.Unlock()

	busy, ok := s.contentTypeBusy[*ct]
	if !ok {
		busy = false
		s.contentTypeBusy[*ct] = busy
	}
	return busy
}

func (s *SubscriptionWatcher) setBusy(ct *ContentType) {
	s.muBusy.Lock()
	defer s.muBusy.Unlock()

	s.contentTypeBusy[*ct] = true
}

func (s *SubscriptionWatcher) unsetBusy(ct *ContentType) {
	s.muBusy.Lock()
	defer s.muBusy.Unlock()

	s.contentTypeBusy[*ct] = false
}

// Run implements the Watcher interface.
func (s *SubscriptionWatcher) Run(ctx context.Context) chan Resource {
	out := make(chan Resource)

	for i := 0; i < contentTypeCount; i++ {
		go s.fetcher(ctx, out)
	}
	go s.generator(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			default:
			}
		}
	}()

	return out
}

// Generator .
func (s *SubscriptionWatcher) generator(ctx context.Context) {
	tickerDur := time.Duration(s.config.TickerIntervalSeconds) * time.Second
	ticker := time.NewTicker(tickerDur)
	defer ticker.Stop()

	go s.generateResources(ctx, time.Now())
	for {
		select {
		default:
			time.Sleep(500 * time.Millisecond)
		case <-ctx.Done():
			close(s.queue)
			return
		case t := <-ticker.C:
			go s.generateResources(ctx, t)
		}
	}
}

func (s *SubscriptionWatcher) generateResources(ctx context.Context, t time.Time) {
	resource := Resource{}

	subscriptions, err := s.client.Subscription.List(ctx)
	if err != nil {
		// TODO: could be a good idea to put the errors
		// TODO: unrelated to a specific contentType audit query
		// TODO: on the SubscriptionWatcher struct.
		// TODO: We would also need to return a separate channel in Run
		// TODO: for sending status/errors to the caller, aside from
		// TODO: the resource channel.
		resource.AddError(err)
		s.sendResourceOrSkip(ctx, resource)
		return
	}

	for _, sub := range subscriptions {

		ct, err := GetContentType(sub.ContentType)
		if err != nil {
			resource.AddError(err)
			s.sendResourceOrSkip(ctx, resource)
			continue
		}
		if s.isBusy(ct) {
			continue
		}
		resource.SetRequest(ct, t)
		s.sendResourceOrSkip(ctx, resource)
	}
}

// Fetcher .
func (s *SubscriptionWatcher) fetcher(ctx context.Context, out chan Resource) {
	for r := range s.queue {
		s.setBusy(r.Request.ContentType)

		lastRequestTime := s.getLastRequestTime(r.Request.ContentType)
		lastContentCreated := s.getLastContentCreated(r.Request.ContentType)

		s.client.logger.Debug(fmt.Sprintf("[%s] lastRequestTime: %s", r.Request.ContentType, lastRequestTime.String()))
		s.client.logger.Debug(fmt.Sprintf("[%s] lastContentCreated: %s", r.Request.ContentType, lastContentCreated.String()))

		start := lastRequestTime
		end := r.Request.RequestTime
		delta := start.Sub(r.Request.RequestTime)
		switch {
		case start.IsZero(), delta < time.Minute:
			lookBehind := time.Duration(s.config.LookBehindMinutes) * time.Minute
			start = r.Request.RequestTime.Add(-(lookBehind))
		case delta > intervalOneDay:
			start = r.Request.RequestTime.Add(-(intervalOneDay))
		}

		s.client.logger.Debug(fmt.Sprintf("[%s] request.RequestTime: %s", r.Request.ContentType, r.Request.RequestTime.String()))
		s.client.logger.Debug(fmt.Sprintf("[%s] fetcher.start: %s", r.Request.ContentType, start.String()))
		s.client.logger.Debug(fmt.Sprintf("[%s] fetcher.end: %s", r.Request.ContentType, end.String()))

		content, err := s.client.Content.List(ctx, r.Request.ContentType, start, end)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				r.AddError(err)
				out <- r
				s.unsetBusy(r.Request.ContentType)
			}
			continue
		}
		s.setLastRequestTime(r.Request.ContentType, r.Request.RequestTime)

		var records []AuditRecord
		for _, c := range content {
			created, err := time.ParseInLocation(CreatedDatetimeFormat, c.ContentCreated, time.Local)
			if err != nil {
				r.AddError(err)
				continue
			}
			s.client.logger.Debug(fmt.Sprintf("[%s] created: %s", r.Request.ContentType, created.String()))
			if !created.After(lastContentCreated) {
				s.client.logger.Debug(fmt.Sprintf("[%s] created skipped", r.Request.ContentType))
				continue
			}
			s.setLastContentCreated(r.Request.ContentType, created)

			s.client.logger.Debug(fmt.Sprintf("[%s] created fetching..", r.Request.ContentType))
			audits, err := s.client.Audit.List(ctx, c.ContentID)
			if err != nil {
				r.AddError(err)
				continue
			}
			records = append(records, audits...)
		}
		select {
		case <-ctx.Done():
			return
		default:
			r.SetResponse(records)
			out <- r
			s.unsetBusy(r.Request.ContentType)
		}
	}
}
