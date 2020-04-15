package office365

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Watcher is an interface used by Watch for generating a stream of records.
type Watcher interface {
	Run(context.Context) chan ResourceAudits
}

// SubscriptionWatcher implements the Watcher interface.
// It fetches current subscriptions, then queries content available for a given interval
// and proceed to query audit records.
type SubscriptionWatcher struct {
	client *Client
	config SubscriptionWatcherConfig
	logger *logrus.Logger

	State
	Handler ResourceHandler
}

// SubscriptionWatcherConfig .
type SubscriptionWatcherConfig struct {
	LookBehindMinutes     int
	TickerIntervalSeconds int
}

// NewSubscriptionWatcher returns a new watcher that uses the provided client
// for querying the API.
func NewSubscriptionWatcher(client *Client, conf SubscriptionWatcherConfig, s State, h ResourceHandler, l *logrus.Logger) (*SubscriptionWatcher, error) {
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
		logger: l,

		State:   s,
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
		s.logger.WithField("content-type", ct.String()).Info("starting worker")
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

		fetch := func(t time.Time) {
			subCh := s.fetchSubscriptions(ctx, done, t)
			for sub := range subCh {
				ctLogger := s.logger.WithField("content-type", sub.ContentType.String())
				workerCh, ok := workers[*sub.ContentType]
				if !ok {
					ctLogger.Error("no worker registered for content-type")
					continue
				}
				select {
				default:
					ctLogger.Warn("worker is busy, skipping")
				case workerCh <- sub:
				}
			}
		}

		fetch(time.Now())
		for {
			select {
			case <-done:
				for ct, workerCh := range workers {
					s.logger.WithField("content-type", ct.String()).Info("closing worker")
					close(workerCh)
				}
				return
			case t := <-ticker.C:
				fetch(t)
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

	return s.Handler.Handle(out)
}

func (s *SubscriptionWatcher) fetchSubscriptions(ctx context.Context, done chan struct{}, t time.Time) chan ResourceSubscription {
	var wg sync.WaitGroup
	out := make(chan ResourceSubscription)

	output := func() {
		defer wg.Done()

		_, subscriptions, err := s.client.Subscription.List(ctx)
		if err != nil {
			subscriptions = []Subscription{}
			if !errors.Is(err, context.Canceled) {
				s.logger.Errorf("fetching subscriptions: %s", err)
			}
		}
		for _, sub := range subscriptions {
			ct, err := GetContentType(sub.ContentType)
			if err != nil {
				s.logger.Errorf("mapping contentType: %s", err)
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

		ctLogger := s.logger.WithField("content-type", sub.ContentType.String())

		end := sub.RequestTime
		ctLogger.Debugf("request.RequestTime: %s", sub.RequestTime.String())

		for {
			lastRequestTime := s.getLastRequestTime(sub.ContentType)
			ctLogger.Debugf("lastRequestTime: %s", lastRequestTime.String())

			start := lastRequestTime
			start, end = s.getTimeWindow(sub.RequestTime, start, end)

			ctLogger.Debugf("fetcher.start: %s", start.String())
			ctLogger.Debugf("fetcher.end: %s", end.String())

			_, content, err := s.client.Content.List(ctx, sub.ContentType, start, end)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					ctLogger.Errorf("could not fetch content: %s", err)
				}
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
			ctLogger := s.logger.WithField("content-type", res.ContentType.String())

			lastContentCreated := s.getLastContentCreated(res.ContentType)
			ctLogger.Debugf("lastContentCreated: %s", lastContentCreated.String())

			created, err := time.ParseInLocation(CreatedDatetimeFormat, res.Content.ContentCreated, time.Local)
			if err != nil {
				ctLogger.Errorf("could not parse ContentCreated: %s", err)
				continue
			}
			ctLogger.Debugf("content found: %s", created.String())
			if !created.After(lastContentCreated) {
				ctLogger.Debugf("content skipped: last[%s] GT current[%s]", lastContentCreated.String(), created.String())
				continue
			}
			s.setLastContentCreated(res.ContentType, created)

			ctLogger.Debugf("content fetching..")
			_, audits, err := s.client.Audit.List(ctx, res.Content.ContentID)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					ctLogger.Errorf("could not fetch audits: %s", err)
				}
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