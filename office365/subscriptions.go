package office365

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// SubscriptionService .
type SubscriptionService service

// List returns the list of subscriptions.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#list-current-subscriptions
//
// List current subscriptions
// This operation returns a collection of the current subscriptions together with the associated webhooks.
func (s *SubscriptionService) List(ctx context.Context) ([]Subscription, error) {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)

	req, err := s.client.newRequest("GET", "subscriptions/list", params.Values, nil)
	if err != nil {
		return nil, err
	}

	var out []Subscription
	_, err = s.client.do(ctx, req, &out)
	return out, err
}

// Start will start a subscription for the specified content type.
// A payload can optionnaly be provided to enable a webhook
// that will send notifications periodically about available content.
// See below webhgook section for details.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#start-a-subscription
//
// This operation starts a subscription to the specified content type. If a subscription to the specified content type already exists, this operation is used to:
// - Update the properties of an active webhook.
// - Enable a webhook that was disabled because of excessive failed notifications.
// - Re-enable an expired webhook by specifying a later or null expiration date.
// - Remove a webhook.
//
// Webhook validation
//
// When the /start operation is called and a webhook is specified, we will send a validation notification
// to the specified webhook address to validate that an active listener can accept and process notifications.
//
// If we do not receive an HTTP 200 OK response, the subscription will not be created.
// Or, if /start is being called to add a webhook to an existing subscription and a response of HTTP 200 OK
// is not received, the webhook will not be added and the subscription will remain unchanged.
func (s *SubscriptionService) Start(ctx context.Context, ct *ContentType, webhook *Webhook) (*Subscription, error) {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)
	if err := params.AddContentType(ct); err != nil {
		return nil, err
	}

	var payload io.Reader
	if webhook != nil {
		data, err := json.Marshal(webhook)
		if err != nil {
			return nil, err
		}
		payload = bytes.NewBuffer(data)
	}

	req, err := s.client.newRequest("POST", "subscriptions/start", params.Values, payload)
	if err != nil {
		return nil, err
	}

	var out *Subscription
	_, err = s.client.do(ctx, req, &out)
	return out, err
}

// Stop stops a subscription for the provided ContentType.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#stop-a-subscription
//
// This operation stops a subscription to the specified content type.
// When a subscription is stopped, you will no longer receive notifications and you will not be able to retrieve available content.
// If the subscription is later restarted, you will have access to new content from that point forward.
// You will not be able to retrieve content that was available between the time the subscription was stopped and restarted.
func (s *SubscriptionService) Stop(ctx context.Context, ct *ContentType) error {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)
	if err := params.AddContentType(ct); err != nil {
		return err
	}

	req, err := s.client.newRequest("POST", "subscriptions/stop", params.Values, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// Watch is used as a dynamic way for fetching events.
// It will poll the current subscriptions for available content
// at regular intervals and returns a channel for consuming returned events.
func (s *SubscriptionService) Watch(ctx context.Context, conf SubscriptionWatcherConfig) (<-chan Resource, error) {
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

		content, err := s.client.Content.List(ctx, r.Request.ContentType, start, end)
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
				audits, err := s.client.Audit.List(ctx, c.ContentID)
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
