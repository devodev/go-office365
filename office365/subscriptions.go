package office365

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"sync"
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
// The context passed will ensure the channel is closed and any underlying
// API queries are notified upon cancellation.
func (s *SubscriptionService) Watch(ctx context.Context, conf SubscriptionWatcherConfig, state State) (<-chan Resource, error) {
	watcher, err := NewSubscriptionWatcher(s.client, conf, state)
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

		fmt.Printf("DEBUG: [%s] lastRequestTime: %s\n", r.Request.ContentType, lastRequestTime.String())
		fmt.Printf("DEBUG: [%s] lastContentCreated: %s\n", r.Request.ContentType, lastContentCreated.String())

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

		fmt.Printf("DEBUG: [%s] request.RequestTime: %s\n", r.Request.ContentType, r.Request.RequestTime.String())
		fmt.Printf("DEBUG: [%s] fetcher.start: %s\n", r.Request.ContentType, start.String())
		fmt.Printf("DEBUG: [%s] fetcher.end: %s\n", r.Request.ContentType, end.String())

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
			fmt.Printf("DEBUG: [%s] created: %s\n", r.Request.ContentType, created.String())

			if !created.After(lastContentCreated) {
				fmt.Printf("DEBUG: [%s] created skipped\n", r.Request.ContentType)
				continue
			}
			s.setLastContentCreated(r.Request.ContentType, created)

			fmt.Printf("DEBUG: [%s] created fetching..\n", r.Request.ContentType)
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

// State is an interface for storing and retrieving Watcher state.
type State interface {
	setLastContentCreated(*ContentType, time.Time)
	getLastContentCreated(*ContentType) time.Time
	setLastRequestTime(*ContentType, time.Time)
	getLastRequestTime(*ContentType) time.Time
}

// MemoryState is an in-memory State interface implementation.
type MemoryState struct {
	muCreated          *sync.RWMutex
	lastContentCreated map[ContentType]time.Time
	muRequest          *sync.RWMutex
	lastRequestTime    map[ContentType]time.Time
}

// NewMemoryState returns a new MemoryState.
func NewMemoryState() *MemoryState {
	return &MemoryState{
		muCreated:          &sync.RWMutex{},
		lastContentCreated: make(map[ContentType]time.Time),
		muRequest:          &sync.RWMutex{},
		lastRequestTime:    make(map[ContentType]time.Time),
	}
}

func (m *MemoryState) setLastContentCreated(ct *ContentType, t time.Time) {
	m.muCreated.Lock()
	defer m.muCreated.Unlock()

	last, ok := m.lastContentCreated[*ct]
	if !ok || last.Before(t) {
		m.lastContentCreated[*ct] = t
	}
}

func (m *MemoryState) getLastContentCreated(ct *ContentType) time.Time {
	m.muCreated.RLock()
	defer m.muCreated.RUnlock()

	t, ok := m.lastContentCreated[*ct]
	if !ok {
		return time.Time{}
	}
	return t
}

func (m *MemoryState) setLastRequestTime(ct *ContentType, t time.Time) {
	m.muRequest.Lock()
	defer m.muRequest.Unlock()

	last, ok := m.lastRequestTime[*ct]
	if !ok || last.Before(t) {
		m.lastRequestTime[*ct] = t
	}
}

func (m *MemoryState) getLastRequestTime(ct *ContentType) time.Time {
	m.muRequest.RLock()
	defer m.muRequest.RUnlock()

	t, ok := m.lastRequestTime[*ct]
	if !ok {
		return time.Time{}
	}
	return t
}

// GOBState is an in-memory State interface implementation, but
// also provides Read and Write methods for serializing/deserializing
// on io.Reader/io.Writer.
// It uses the encoding/gob package.
type GOBState struct {
	*MemoryState
}

// NewGOBState returns a new GOBState.
func NewGOBState() *GOBState {
	return &GOBState{NewMemoryState()}
}

func (g *GOBState) createBlob() *GOBStateBlob {
	g.muCreated.RLock()
	g.muRequest.RLock()
	defer g.muCreated.RUnlock()
	defer g.muRequest.RUnlock()

	return &GOBStateBlob{
		LastContentCreated: g.lastContentCreated,
		LastRequestTime:    g.lastRequestTime,
	}
}

func (g *GOBState) setFromBlob(b *GOBStateBlob) {
	g.muCreated.Lock()
	g.muRequest.Lock()
	defer g.muCreated.Unlock()
	defer g.muRequest.Unlock()

	g.lastContentCreated = b.LastContentCreated
	g.lastRequestTime = b.LastRequestTime
}

// Read will deserialize from a reader and populate its internal state.
func (g *GOBState) Read(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	var blob GOBStateBlob
	if err := decoder.Decode(&blob); err != nil {
		return err
	}
	g.setFromBlob(&blob)
	return nil
}

// Write will serialize its internal state and write to a writer.
func (g *GOBState) Write(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	blob := g.createBlob()
	if err := encoder.Encode(&blob); err != nil {
		return err
	}
	return nil
}

// GOBStateBlob is used to serialize/deserialize MemoryState
// internal state.
type GOBStateBlob struct {
	LastContentCreated map[ContentType]time.Time
	LastRequestTime    map[ContentType]time.Time
}

// Resource .
type Resource struct {
	Request  ResourceRequest
	Response ResourceResponse
	Errors   []error
}

// AddError .
func (r *Resource) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

// SetRequest .
func (r *Resource) SetRequest(ct *ContentType, t time.Time) {
	r.Request = ResourceRequest{
		ContentType: ct,
		RequestTime: t,
	}
}

// SetResponse .
func (r *Resource) SetResponse(records []AuditRecord) {
	r.Response = ResourceResponse{records}
}

// ResourceRequest .
type ResourceRequest struct {
	ContentType *ContentType
	RequestTime time.Time
}

// ResourceResponse .
type ResourceResponse struct {
	Records []AuditRecord
}

// ResourceHandler is an interface for handling streamed resources.
type ResourceHandler interface {
	Handle(<-chan Resource)
}

// Printer implements the ResourceHandler interface.
// It prints a human readable formatted resource on the
// provided writer.
type Printer struct {
	writer io.Writer
}

// NewPrinter returns a printer using the provided writer.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w}
}

// Handle .
func (h Printer) Handle(in <-chan Resource) {
	for r := range in {
		for idx, e := range r.Errors {
			fmt.Fprintf(h.writer, "[%s] Error%d: %s", r.Request.ContentType, idx, e.Error())
		}
		for _, a := range r.Response.Records {
			auditStr, err := json.Marshal(a)
			if err != nil {
				fmt.Fprintf(h.writer, "error marshalling audit: %s\n", err)
				continue
			}
			var out bytes.Buffer
			json.Indent(&out, auditStr, "", "\t")
			fmt.Fprintf(h.writer, "[%s]\n%s\n", r.Request.ContentType, out.String())
		}
	}
}
