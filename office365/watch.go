package office365

import (
	"context"
	"fmt"
	"time"
)

// WatchService .
type WatchService service

// Watch is used as a dynamic way for fetching events.
// It will poll the current subscriptions for available content
// at regular intervals and returns a channel for consuming returned events.
func (s *WatchService) Watch(ctx context.Context, conf SubscriptionWatcherConfig) (<-chan Resource, error) {
	watcher, err := NewSubscriptionWatcher(s.client, conf)
	if err != nil {
		return nil, err
	}

	resourceChan := watcher.Run(ctx)

	return resourceChan, nil
}

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

	queue chan Resource
}

// SubscriptionWatcherConfig .
type SubscriptionWatcherConfig struct {
	FetcherCount           int
	LookBehindMinutes      int
	FetcherIntervalSeconds int
	TickerIntervalSeconds  int
}

// NewSubscriptionWatcher returns a new watcher that uses the provided client
// for querying the API.
func NewSubscriptionWatcher(client *Client, conf SubscriptionWatcherConfig) (*SubscriptionWatcher, error) {
	if conf.FetcherCount <= 0 {
		return nil, fmt.Errorf("fetcherCount must be greater than 0")
	}

	lookBehindDur := time.Duration(conf.LookBehindMinutes) * time.Minute
	if lookBehindDur <= 0 {
		return nil, fmt.Errorf("lookBehindMinutes must be greater than 0")
	}
	if lookBehindDur > 24*time.Hour {
		return nil, fmt.Errorf("lookBehindMinutes must be less than 24 hours")
	}

	fetcherIntervalDur := time.Duration(conf.FetcherIntervalSeconds) * time.Second
	if fetcherIntervalDur <= 0 {
		return nil, fmt.Errorf("fetcherIntervalSeconds must be greater than 0")
	}
	if fetcherIntervalDur > 24*time.Hour {
		return nil, fmt.Errorf("fetcherIntervalSeconds must be less than 24 hours")
	}

	tickerIntervalDur := time.Duration(conf.TickerIntervalSeconds) * time.Second
	if tickerIntervalDur <= 0 {
		return nil, fmt.Errorf("tickerIntervalSeconds must be greater than 0")
	}
	if tickerIntervalDur > 24*time.Hour {
		return nil, fmt.Errorf("tickerIntervalSeconds must be less than 24 hours")
	}

	watcher := &SubscriptionWatcher{
		client: client,
		config: conf,

		queue: make(chan Resource),
	}
	return watcher, nil
}

// Run implements the Watcher interface.
func (s SubscriptionWatcher) Run(ctx context.Context) chan Resource {
	out := make(chan Resource)

	for i := 0; i < s.config.FetcherCount; i++ {
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
func (s SubscriptionWatcher) generator(ctx context.Context) {
	fetcherDur := time.Duration(s.config.FetcherIntervalSeconds) * time.Second
	tickerDur := time.Duration(s.config.TickerIntervalSeconds) * time.Second
	ticker := time.NewTicker(tickerDur)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(s.queue)
			return
		case t := <-ticker.C:
			go func() {
				resource := Resource{}

				subscriptions, err := s.client.Subscription.List(ctx)
				if err != nil {
					resource.AddError(err)
					s.queue <- resource
					return
				}

				startTime := t.Add(-(fetcherDur))
				endTime := t

				for _, sub := range subscriptions {
					ct, err := GetContentType(sub.ContentType)
					if err != nil {
						resource.AddError(err)
						s.queue <- resource
						continue
					}
					resource.SetRequest(ct, startTime, endTime)
					s.queue <- resource
				}
			}()
		}
	}
}

// Fetcher .
func (s SubscriptionWatcher) fetcher(ctx context.Context, out chan Resource) {
	for r := range s.queue {
		lookBehind := time.Duration(s.config.LookBehindMinutes) * time.Minute
		start := r.Request.EndTime.Add(-(lookBehind))
		end := r.Request.EndTime

		content, err := s.client.Content.Content(ctx, r.Request.ContentType, start, end)
		if err != nil {
			r.AddError(err)
			out <- r
			continue
		}

		fmt.Printf("DEBUG: [%s] fetcher.start: %s\n", r.Request.ContentType, start.String())
		fmt.Printf("DEBUG: [%s] fetcher.end: %s\n", r.Request.ContentType, end.String())

		fmt.Printf("DEBUG: [%s] request.startTime: %s\n", r.Request.ContentType, r.Request.StartTime.String())
		fmt.Printf("DEBUG: [%s] request.EndTime: %s\n", r.Request.ContentType, r.Request.EndTime.String())

		var records []AuditRecord
		for _, c := range content {
			created, err := time.ParseInLocation(CreatedDatetimeFormat, c.ContentCreated, time.Local)
			if err != nil {
				r.AddError(err)
				continue
			}
			fmt.Printf("DEBUG: [%s] created: %s\n", r.Request.ContentType, created.String())

			createdAfterOrEqual := created.After(r.Request.StartTime) || created.Equal(r.Request.StartTime)
			if createdAfterOrEqual && created.Before(r.Request.EndTime) {
				audits, err := s.client.Audit.Audit(ctx, c.ContentID)
				if err != nil {
					r.AddError(err)
					continue
				}
				records = append(records, audits...)
			}
		}
		r.SetResponse(records)
		out <- r
	}
}
