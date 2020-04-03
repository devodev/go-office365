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
	}
	return watcher, nil
}

// Run implements the Watcher interface.
func (s *SubscriptionWatcher) Run(ctx context.Context) chan ResourceAudits {
	done := make(chan struct{})
	out := make(chan ResourceAudits)

	go func(d chan struct{}, o chan ResourceAudits) {
		defer close(o)

		tickerDur := time.Duration(s.config.TickerIntervalSeconds) * time.Second
		ticker := time.NewTicker(tickerDur)
		defer ticker.Stop()

		for {
			select {
			case <-d:
				return
			case t := <-ticker.C:
				subCh := s.fetchSubscriptions(ctx, d)
				contentCh := s.fetchContent(ctx, d, subCh, t)
				auditCh := s.fetchAudits(ctx, d, contentCh)

				for a := range auditCh {
					o <- a
				}
			}
		}
	}(done, out)

	go func() {
		select {
		case <-ctx.Done():
			close(done)
			return
		}
	}()

	return out
}

func (s *SubscriptionWatcher) fetchSubscriptions(ctx context.Context, done chan struct{}) chan ResourceSubscription {
	var wg sync.WaitGroup
	out := make(chan ResourceSubscription)

	output := func(ch chan ResourceSubscription) {
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
			case ch <- ResourceSubscription{ct, sub}:
			}
		}
	}

	wg.Add(1)
	go output(out)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (s *SubscriptionWatcher) fetchContent(ctx context.Context, done chan struct{}, subCh chan ResourceSubscription, t time.Time) chan ResourceContent {
	var wg sync.WaitGroup
	out := make(chan ResourceContent)

	output := func(ch chan ResourceContent) {
		defer wg.Done()

	Outer:
		for sub := range subCh {
			end := t
			s.client.logger.Debug(fmt.Sprintf("[%s] request.RequestTime: %s", sub.ContentType, t.String()))

			for {
				lastRequestTime := s.getLastRequestTime(sub.ContentType)
				s.client.logger.Debug(fmt.Sprintf("[%s] lastRequestTime: %s", sub.ContentType, lastRequestTime.String()))

				start := lastRequestTime
				start, end = s.getTimeWindow(t, start, end)

				s.client.logger.Debug(fmt.Sprintf("[%s] fetcher.start: %s", sub.ContentType, start.String()))
				s.client.logger.Debug(fmt.Sprintf("[%s] fetcher.end: %s", sub.ContentType, end.String()))

				content, err := s.client.Content.List(ctx, sub.ContentType, start, end)
				if err != nil {
					s.client.logger.Printf("error fetching content for %s: %s", sub.ContentType.String(), err)
					continue Outer
				}
				for _, c := range content {
					select {
					case <-done:
						return
					case ch <- ResourceContent{sub.ContentType, t, c}:
					}
				}
				s.setLastRequestTime(sub.ContentType, end)
				if !end.Before(t) {
					break
				}
			}
		}
	}

	wg.Add(1)
	go output(out)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (s *SubscriptionWatcher) fetchAudits(ctx context.Context, done chan struct{}, contentCh chan ResourceContent) chan ResourceAudits {
	var wg sync.WaitGroup
	out := make(chan ResourceAudits)

	output := func(ch chan ResourceAudits) {
		defer wg.Done()

		for res := range contentCh {
			lastContentCreated := s.getLastContentCreated(res.ContentType)
			s.client.logger.Debug(fmt.Sprintf("[%s] lastContentCreated: %s", res.ContentType, lastContentCreated.String()))

			created, err := time.ParseInLocation(CreatedDatetimeFormat, res.Content.ContentCreated, time.Local)
			if err != nil {
				s.client.logger.Printf("error fetching audit for %s: %s", res.ContentType, err)
				continue
			}
			s.client.logger.Debug(fmt.Sprintf("[%s] created: %s", res.ContentType, created.String()))
			if !created.After(lastContentCreated) {
				s.client.logger.Debug(fmt.Sprintf("[%s] created skipped", res.ContentType))
				continue
			}
			s.setLastContentCreated(res.ContentType, created)

			s.client.logger.Debug(fmt.Sprintf("[%s] created fetching..", res.ContentType))
			audits, err := s.client.Audit.List(ctx, res.Content.ContentID)
			if err != nil {
				s.client.logger.Printf("error fetching audits for %s: %s", res.ContentType, err)
				continue
			}
			for _, a := range audits {
				select {
				case <-done:
					return
				case ch <- ResourceAudits{res.ContentType, res.RequestTime, a}:
				}
			}
		}
	}

	wg.Add(1)
	go output(out)

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
