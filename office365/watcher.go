package office365

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Watcher is an interface used by Watch for generating a stream of records.
type Watcher interface {
	Run(context.Context) chan ResourceAudits
}

// SubscriptionWatcher implements the Watcher interface.
// It fecthes current subscriptions, then queries content available for a given interval
// and proceed to query audit records.
type SubscriptionWatcher struct {
	client *Client
	config SubscriptionWatcherConfig

	State
	Handler ResourceHandler
}

// SubscriptionWatcherConfig .
type SubscriptionWatcherConfig struct {
	LookBehindMinutes      int
	TickerIntervalSeconds  int
	RefreshIntervalMinutes int
}

// NewSubscriptionWatcher returns a new watcher that uses the provided client
// for querying the API.
func NewSubscriptionWatcher(client *Client, conf SubscriptionWatcherConfig, s State, h ResourceHandler) (*SubscriptionWatcher, error) {
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

		Handler: h,
	}
	return watcher, nil
}

// Run implements the Watcher interface.
func (s *SubscriptionWatcher) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	done := make(chan struct{})
	out := make(chan ResourceAudits)

	// setup worker pool
	// workers receive jobs and send results to output channel
	workers := make(map[ContentType]chan ResourceSubscription)
	contentTypes := GetContentTypes()

	wg.Add(len(contentTypes))
	for _, ct := range contentTypes {
		ch := make(chan ResourceSubscription, 1)
		workers[ct] = ch

		go func() {
			defer wg.Done()
			for res := range ch {
				contentCh := s.fetchContent(ctx, done, res)
				auditCh := s.fetchAudits(ctx, done, contentCh)

				for a := range auditCh {
					out <- a
				}
			}
		}()
	}

	// this goroutine is responsible for closing output channel
	go func() {
		wg.Wait()
		close(out)
	}()

	// setup ticker that will periodically fetch subscriptions
	// and create jobs for workers.
	// this goroutine is responsible for closing worker channels
	go func() {
		tickerDur := time.Duration(s.config.TickerIntervalSeconds) * time.Second
		ticker := time.NewTicker(tickerDur)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				for ct, workerCh := range workers {
					s.client.logger.Printf("closing worker for %q", ct.String())
					close(workerCh)
				}
				return
			case t := <-ticker.C:
				subCh := s.fetchSubscriptions(ctx, done, t)
				for sub := range subCh {
					workerCh, ok := workers[*sub.ContentType]
					if !ok {
						s.client.logger.Printf("no worker available for %q", sub.ContentType.String())
						continue
					}
					select {
					default:
						s.client.logger.Printf("skipping %q because worker is busy", sub.ContentType.String())
					case workerCh <- sub:
					}
				}
			}
		}
	}()

	// this goroutine is responsible for notifying
	// everyone that we want to exit
	go func() {
		select {
		case <-ctx.Done():
			close(done)
			return
		}
	}()

	return s.Handler.Handle(out, s.client.logger)
}

func (s *SubscriptionWatcher) fetchSubscriptions(ctx context.Context, done chan struct{}, t time.Time) chan ResourceSubscription {
	var wg sync.WaitGroup
	out := make(chan ResourceSubscription)

	output := func() {
		defer wg.Done()

		subscriptions, err := s.client.Subscription.List(ctx)
		if err != nil {
			subscriptions = []Subscription{}
			s.client.logger.Printf("error fetching subscriptions: %s", err)
		}
		for _, sub := range subscriptions {
			ct, err := GetContentType(sub.ContentType)
			if err != nil {
				s.client.logger.Printf("error mapping contentType: %s", err)
				continue
			}
			select {
			case <-done:
				return
			case out <- ResourceSubscription{ct, t, sub}:
			}
		}
	}

	wg.Add(1)
	go output()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (s *SubscriptionWatcher) fetchContent(ctx context.Context, done chan struct{}, res ResourceSubscription) chan ResourceContent {
	var wg sync.WaitGroup
	out := make(chan ResourceContent)

	output := func(sub ResourceSubscription) {
		defer wg.Done()

		end := sub.RequestTime
		s.client.logger.Printf("[%s] request.RequestTime: %s", sub.ContentType, sub.RequestTime.String())

		for {
			lastRequestTime := s.getLastRequestTime(sub.ContentType)
			s.client.logger.Printf("[%s] lastRequestTime: %s", sub.ContentType, lastRequestTime.String())

			start := lastRequestTime
			start, end = s.getTimeWindow(sub.RequestTime, start, end)

			s.client.logger.Printf("[%s] fetcher.start: %s", sub.ContentType, start.String())
			s.client.logger.Printf("[%s] fetcher.end: %s", sub.ContentType, end.String())

			content, err := s.client.Content.List(ctx, sub.ContentType, start, end)
			if err != nil {
				s.client.logger.Printf("could not fetch content for %s: %s", sub.ContentType, err)
				return
			}
			for _, c := range content {
				select {
				case <-done:
					return
				case out <- ResourceContent{sub.ContentType, sub.RequestTime, c}:
				}
			}
			s.setLastRequestTime(sub.ContentType, end)
			if !end.Before(sub.RequestTime) {
				break
			}
		}
	}

	wg.Add(1)
	go output(res)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (s *SubscriptionWatcher) fetchAudits(ctx context.Context, done chan struct{}, contentCh chan ResourceContent) chan ResourceAudits {
	var wg sync.WaitGroup
	out := make(chan ResourceAudits)

	output := func(ch <-chan ResourceContent) {
		defer wg.Done()

		for res := range ch {
			lastContentCreated := s.getLastContentCreated(res.ContentType)
			s.client.logger.Printf("[%s] lastContentCreated: %s", res.ContentType, lastContentCreated.String())

			created, err := time.ParseInLocation(CreatedDatetimeFormat, res.Content.ContentCreated, time.Local)
			if err != nil {
				s.client.logger.Printf("could not parse ContentCreated for %s: %s", res.ContentType, err)
				continue
			}
			s.client.logger.Printf("[%s] content found: %s", res.ContentType, created.String())
			if !created.After(lastContentCreated) {
				s.client.logger.Printf("[%s] content skipped: last:%s >= current:%s",
					res.ContentType, lastContentCreated.String(), created.String())
				continue
			}
			s.setLastContentCreated(res.ContentType, created)

			s.client.logger.Printf("[%s] content fetching..", res.ContentType)
			audits, err := s.client.Audit.List(ctx, res.Content.ContentID)
			if err != nil {
				s.client.logger.Printf("could not fetch audits for %s: %s", res.ContentType, err)
				continue
			}
			for _, a := range audits {
				select {
				case <-done:
					return
				case out <- ResourceAudits{res.ContentType, res.RequestTime, a}:
				}
			}
		}
	}

	wg.Add(1)
	go output(contentCh)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
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

// ResourceSubscription .
type ResourceSubscription struct {
	ContentType  *ContentType
	RequestTime  time.Time
	Subscription Subscription
}

// ResourceContent .
type ResourceContent struct {
	ContentType *ContentType
	RequestTime time.Time
	Content     Content
}

// ResourceAudits .
type ResourceAudits struct {
	ContentType *ContentType
	RequestTime time.Time
	AuditRecord AuditRecord
}
