package office365

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
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

	subscriptions, err := client.Subscription.List(context.Background())
	if err != nil {
		t.Errorf("error occured running Subscriptions.List: %v", err)
	}

	want := []Subscription{
		Subscription{
			ContentType: "test",
			Status:      "test",
			Webhook:     Webhook{},
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
			if webhook.Address == "" {
				t.Errorf("webhook.address is required")
			}
		}

		response := &Subscription{
			ContentType: contentType,
			Status:      "enabled",
		}
		if webhook != nil {
			response.Webhook = Webhook{
				Status:     "enabled",
				Address:    webhook.Address,
				AuthID:     webhook.AuthID,
				Expiration: webhook.Expiration,
			}
		}

		json.NewEncoder(w).Encode(response)
	})

	cases := []struct {
		Request     *Webhook
		ContentType ContentType
		Want        *Subscription
	}{
		{
			Request: &Webhook{
				Address:    "test-address",
				AuthID:     "test-authid",
				Expiration: "",
			},
			ContentType: AuditAzureActiveDirectory,
			Want: &Subscription{
				ContentType: AuditAzureActiveDirectory.String(),
				Status:      "enabled",
				Webhook: Webhook{
					Status:     "enabled",
					Address:    "test-address",
					AuthID:     "test-authid",
					Expiration: "",
				},
			},
		},
		{
			Request:     nil,
			ContentType: AuditAzureActiveDirectory,
			Want: &Subscription{
				ContentType: AuditAzureActiveDirectory.String(),
				Status:      "enabled",
				Webhook:     Webhook{},
			},
		},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d.", idx), func(t *testing.T) {
			subscriptions, err := client.Subscription.Start(context.Background(), &c.ContentType, c.Request)
			if err != nil {
				t.Errorf("error occured running Subscriptions.Start: %v", err)
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
		ContentType ContentType
	}{
		{ContentType: AuditAzureActiveDirectory},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d.", idx), func(t *testing.T) {
			err := client.Subscription.Stop(context.Background(), &c.ContentType)
			if err != nil {
				t.Errorf("error occured running Subscriptions.Stop: %v", err)
			}
		})
	}
}
