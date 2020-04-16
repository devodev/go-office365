package office365

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

// SubscriptionService .
type SubscriptionService service

// List returns the list of subscriptions.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#list-current-subscriptions
//
// List current subscriptions
// This operation returns a collection of the current subscriptions together with the associated webhooks.
func (s *SubscriptionService) List(ctx context.Context) (*Response, []Subscription, error) {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)

	req, err := s.client.newRequest("GET", "subscriptions/list", params.Values, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []Subscription
	resp, err := s.client.do(ctx, req, &out)
	return resp, out, err
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
func (s *SubscriptionService) Start(ctx context.Context, ct *ContentType, webhook *Webhook) (*Response, *Subscription, error) {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)
	if err := params.AddContentType(ct); err != nil {
		return nil, nil, err
	}

	var payload io.Reader
	if webhook != nil {
		data, err := json.Marshal(webhook)
		if err != nil {
			return nil, nil, err
		}
		payload = bytes.NewBuffer(data)
	}

	req, err := s.client.newRequest("POST", "subscriptions/start", params.Values, payload)
	if err != nil {
		return nil, nil, err
	}

	var out *Subscription
	resp, err := s.client.do(ctx, req, &out)
	return resp, out, err
}

// Stop stops a subscription for the provided ContentType.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#stop-a-subscription
//
// This operation stops a subscription to the specified content type.
// When a subscription is stopped, you will no longer receive notifications and you will not be able to retrieve available content.
// If the subscription is later restarted, you will have access to new content from that point forward.
// You will not be able to retrieve content that was available between the time the subscription was stopped and restarted.
func (s *SubscriptionService) Stop(ctx context.Context, ct *ContentType) (*Response, error) {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)
	if err := params.AddContentType(ct); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest("POST", "subscriptions/stop", params.Values, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.do(ctx, req, nil)
	return resp, err
}

// Subscription represents a response.
type Subscription struct {
	ContentType string   `json:"contentType"`
	Status      string   `json:"status"`
	Webhook     *Webhook `json:"webhook"`
}

// Webhook represents both a response and a request payload.
type Webhook struct {
	Status     string `json:"status,omitempty"`
	Address    string `json:"address"`
	AuthID     string `json:"authId,omitempty"`
	Expiration string `json:"expiration,omitempty"`
}
