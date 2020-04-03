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

	State

	mu            *sync.Mutex
	subscriptions map[ContentType]chan Resource
}

// SubscriptionWatcherConfig .
type SubscriptionWatcherConfig struct {
	LookBehindMinutes      int
	TickerIntervalSeconds  int
	RefreshIntervalMinutes int
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

		State: s,

		mu:            &sync.Mutex{},
		subscriptions: make(map[ContentType]chan Resource),
	}
	return watcher, nil
}

// Run implements the Watcher interface.
func (s *SubscriptionWatcher) Run(ctx context.Context) chan Resource {
	out := make(chan Resource)

	refreshTickerDur := time.Duration(s.config.RefreshIntervalMinutes) * time.Minute
	refreshTicker := time.NewTicker(refreshTickerDur)
	defer refreshTicker.Stop()
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-refreshTicker.C:
			go s.refreshSubscriptions(ctx, out)
		}
	}()

	tickerDur := time.Duration(s.config.TickerIntervalSeconds) * time.Second
	ticker := time.NewTicker(tickerDur)
	defer ticker.Stop()
	go func() {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			go s.generateResources(ctx, t)
		}
	}()

	return out
}

func (s *SubscriptionWatcher) refreshSubscriptions(ctx context.Context, out chan Resource) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ct, ch := range s.subscriptions {
		close(ch)
		delete(s.subscriptions, ct)
	}

	subscriptions, err := s.client.Subscription.List(ctx)
	if err != nil {
		s.client.logger.Printf("error while fetching subscriptions: %s", err)
		return
	}

	for _, sub := range subscriptions {
		ct, err := GetContentType(sub.ContentType)
		if err != nil {
			s.client.logger.Printf("error while mapping contentType: %s", err)
			continue
		}
		inResourceChan := make(chan Resource)
		s.subscriptions[*ct] = inResourceChan

		go s.fetcher(ctx, inResourceChan, out)
	}
}

func (s *SubscriptionWatcher) generateResources(ctx context.Context, t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ct, ch := range s.subscriptions {
		resource := Resource{}
		resource.SetRequest(&ct, t)
		select {
		default:
		case ch <- resource:
		}
	}
}

func (s *SubscriptionWatcher) fetcher(ctx context.Context, in, out chan Resource) {
Outer:
	for r := range in {
		lastContentCreated := s.getLastContentCreated(r.Request.ContentType)
		s.client.logger.Debug(fmt.Sprintf("[%s] lastContentCreated: %s", r.Request.ContentType, lastContentCreated.String()))

		end := r.Request.RequestTime
		s.client.logger.Debug(fmt.Sprintf("[%s] request.RequestTime: %s", r.Request.ContentType, r.Request.RequestTime.String()))

		var finalContent []Content
		for {
			lastRequestTime := s.getLastRequestTime(r.Request.ContentType)
			s.client.logger.Debug(fmt.Sprintf("[%s] lastRequestTime: %s", r.Request.ContentType, lastRequestTime.String()))

			start := lastRequestTime
			start, end = s.getTimeWindow(r.Request.RequestTime, start, end)

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
				}
				continue Outer
			}
			finalContent = append(finalContent, content...)

			s.setLastRequestTime(r.Request.ContentType, end)

			if !end.Before(r.Request.RequestTime) {
				break
			}
		}

		var records []AuditRecord
		for _, c := range finalContent {
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
		}
	}
}

func (s *SubscriptionWatcher) getTimeWindow(requestTime, start, end time.Time) (time.Time, time.Time) {
	if start.Equal(end) {
		end = requestTime
	}

	delta := end.Sub(start)
	lookbehindDelta := time.Duration(s.config.LookBehindMinutes) * time.Minute

	switch {
	case start.IsZero(), start.After(end), delta < lookbehindDelta:
		// base case
		// we move the start behind
		start = end.Add(-(lookbehindDelta))
	case end.Before(requestTime):
		// we have looped, adjust the end
		end.Add(intervalOneDay)
	case delta > intervalOneWeek:
		// cant query API later than one week in the past
		// move the interval window behind
		start = end.Add(-(intervalOneWeek))
		end = start.Add(intervalOneDay)
	case delta > intervalOneDay:
		// cant query API for more than 24 hour interval
		// we move the end behind
		end = start.Add(intervalOneDay)
	}
	if end.After(requestTime) {
		end = requestTime
	}
	return start, end
}
