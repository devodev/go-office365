package office365

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/devodev/go-office365/v0/pkg/office365/schema"
)

func TestList(t *testing.T) {

	client, mux, teardown := stubClient()
	defer teardown()

	url := client.getURL("subscriptions/list", nil)
	mux.HandleFunc(url.Path, func(w http.ResponseWriter, r *http.Request) {
		EnforceMethod(t, r, "GET")
		fmt.Fprint(w, `[
            {
                "contentType": "test",
                "status": "test",
                "webhook": null
            }
        ]`)
	})

	_, subscriptions, err := client.Subscription.List(context.Background())
	if err != nil {
		t.Errorf("error occurred running Subscriptions.List: %v", err)
	}

	want := []Subscription{
		{
			ContentType: String("test"),
			Status:      String("test"),
			Webhook:     nil,
		},
	}
	testDeep(t, subscriptions, want)
}

func TestStart(t *testing.T) {

	client, mux, teardown := stubClient()
	defer teardown()

	url := client.getURL("subscriptions/start", nil)
	mux.HandleFunc(url.Path, func(w http.ResponseWriter, r *http.Request) {
		EnforceMethod(t, r, "POST")
		contentType := EnforceAndReturnContentType(t, r)

		var webhook *Webhook
		if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
			if err != io.EOF {
				t.Errorf("error decoding body: %s", err)
			}
		}
		if webhook != nil {
			if webhook.Address == nil {
				t.Errorf("webhook.address is required")
			}
		}

		response := &Subscription{
			ContentType: &contentType,
			Status:      String("enabled"),
		}
		if webhook != nil {
			response.Webhook = &Webhook{
				Status:     String("enabled"),
				Address:    webhook.Address,
				AuthID:     webhook.AuthID,
				Expiration: webhook.Expiration,
			}
		}

		json.NewEncoder(w).Encode(response)
	})

	cases := []struct {
		Request     *Webhook
		ContentType schema.ContentType
		Want        *Subscription
	}{
		{
			Request: &Webhook{
				Address:    String("test-address"),
				AuthID:     String("test-authid"),
				Expiration: String(""),
			},
			ContentType: schema.AuditAzureActiveDirectory,
			Want: &Subscription{
				ContentType: String(schema.AuditAzureActiveDirectory.String()),
				Status:      String("enabled"),
				Webhook: &Webhook{
					Status:     String("enabled"),
					Address:    String("test-address"),
					AuthID:     String("test-authid"),
					Expiration: String(""),
				},
			},
		},
		{
			Request:     nil,
			ContentType: schema.AuditAzureActiveDirectory,
			Want: &Subscription{
				ContentType: String(schema.AuditAzureActiveDirectory.String()),
				Status:      String("enabled"),
				Webhook:     nil,
			},
		},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d.", idx), func(t *testing.T) {
			_, subscriptions, err := client.Subscription.Start(context.Background(), &c.ContentType, c.Request)
			if err != nil {
				t.Errorf("error occurred running Subscriptions.Start: %v", err)
			}
			testDeep(t, subscriptions, c.Want)
		})
	}
}

func TestStop(t *testing.T) {

	client, mux, teardown := stubClient()
	defer teardown()

	url := client.getURL("subscriptions/stop", nil)
	mux.HandleFunc(url.Path, func(w http.ResponseWriter, r *http.Request) {
		EnforceMethod(t, r, "POST")
		EnforceAndReturnContentType(t, r)
	})

	cases := []struct {
		ContentType schema.ContentType
	}{
		{ContentType: schema.AuditAzureActiveDirectory},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d.", idx), func(t *testing.T) {
			_, err := client.Subscription.Stop(context.Background(), &c.ContentType)
			if err != nil {
				t.Errorf("error occurred running Subscriptions.Stop: %v", err)
			}
		})
	}
}
