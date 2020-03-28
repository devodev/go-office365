package office365

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func stubClient() (*Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := NewClient(nil, "test-tenandID", "")
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url

	return client, mux, server.Close
}

func EnforceMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("got request method: %v but want %v", got, want)
	}
}

func EnforceAndReturnContentType(t *testing.T, r *http.Request) string {
	t.Helper()
	contentType := r.URL.Query().Get("contentType")
	if contentType == "" {
		t.Errorf("contentType queryParam is required")
		return ""
	}
	if !ContentTypeValid(contentType) {
		t.Errorf("unknown contentType: %s", contentType)
		return ""
	}
	return contentType
}

func EnforceAndReturnTime(t *testing.T, r *http.Request, param string) time.Time {
	t.Helper()
	parsed, err := time.Parse(RequestDateFormat, param)
	if err == nil {
		t.Logf("parsed time successfully using RequestDateFormat: %v == [%v] ", param, RequestDateFormat)
		return parsed
	}
	t.Logf("X could not parse time using RequestDateFormat: %v != [%v] ", param, RequestDateFormat)
	parsed, err = time.Parse(RequestDatetimeFormat, param)
	if err == nil {
		t.Logf("parsed time successfully using RequestDatetimeFormat: %v == [%v] ", param, RequestDatetimeFormat)
		return parsed
	}
	t.Logf("X could not parse time using RequestDatetimeFormat: %v != [%v] ", param, RequestDatetimeFormat)
	parsed, err = time.Parse(RequestDatetimeLargeFormat, param)
	if err == nil {
		t.Logf("parsed time successfully using RequestDatetimeLargeFormat: %v == [%v] ", param, RequestDatetimeLargeFormat)
		return parsed
	}
	t.Logf("X could not parse time using RequestDatetimeLargeFormat: %v != [%v] ", param, RequestDatetimeLargeFormat)
	return time.Time{}
}

func testError(t *testing.T, want interface{}, wantError error, gotError error) {
	t.Helper()
	if gotError != nil {
		if wantError != nil && gotError.Error() != wantError.Error() {
			t.Errorf("error occured but different than WantError: %v != %v", gotError, wantError)
		}
		if wantError == nil {
			t.Errorf("error occured but WantError is nil: %v", gotError)
		}
		return
	}
	if want == nil {
		t.Errorf("no error occured but Want is nil")
		return
	}
}

func testDeep(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got\n%v\nbut want\n%v", got, want)
	}
}

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

	subscriptions, err := client.Subscriptions.List(context.Background())
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
			subscriptions, err := client.Subscriptions.Start(context.Background(), &c.ContentType, c.Request)
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
			err := client.Subscriptions.Stop(context.Background(), &c.ContentType)
			if err != nil {
				t.Errorf("error occured running Subscriptions.Stop: %v", err)
			}
		})
	}
}

func TestContent(t *testing.T) {

	client, mux, teardown := stubClient()
	defer teardown()

	now := time.Now()
	nowMinusintervalOneDay := now.Add(-intervalOneDay)
	contentID := "test-contentid"
	contentURI := client.getURL(contentID, nil)
	store := []Content{
		{
			ContentType:       AuditAzureActiveDirectory.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Add(-(intervalOneDay * 6)).Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 1).Format(time.RFC3339),
		},
		{
			ContentType:       AuditAzureActiveDirectory.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Add(-(intervalOneDay * 5)).Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 2).Format(time.RFC3339),
		},
		{
			ContentType:       AuditSharePoint.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Add(-(intervalOneDay * 2)).Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 5).Format(time.RFC3339),
		},
		{
			ContentType:       DLPAll.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 7).Format(time.RFC3339),
		},
		// test next-uri header
		{
			ContentType:       AuditExchange.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Add(-(intervalOneDay * 2)).Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 5).Format(time.RFC3339),
		},
		{
			ContentType:       AuditExchange.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Add(-(intervalOneDay * 2)).Add(time.Minute).Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 5).Add(time.Minute).Format(time.RFC3339),
		},
		{
			ContentType:       AuditExchange.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Add(-(intervalOneDay * 2)).Add(time.Minute * 2).Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 5).Add(time.Minute * 2).Format(time.RFC3339),
		},
		{
			ContentType:       AuditExchange.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Add(-(intervalOneDay * 2)).Add(time.Minute * 3).Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 5).Add(time.Minute * 3).Format(time.RFC3339),
		},
		{
			ContentType:       AuditExchange.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Add(-(intervalOneDay * 2)).Add(time.Minute * 4).Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 5).Add(time.Minute * 4).Format(time.RFC3339),
		},
		{
			ContentType:       AuditExchange.String(),
			ContentID:         contentID,
			ContentURI:        contentURI.String(),
			ContentCreated:    now.Add(-(intervalOneDay * 2)).Add(time.Minute * 5).Format(time.RFC3339),
			ContentExpiration: now.Add(intervalOneDay * 5).Add(time.Minute * 5).Format(time.RFC3339),
		},
	}

	filterStore := func(s *[]Content, contentType string, startTime time.Time, EndTime time.Time) []Content {
		var result []Content
		for _, v := range *s {
			created, _ := time.Parse(time.RFC3339, v.ContentCreated)
			if v.ContentType == contentType {
				if startTime.IsZero() && EndTime.IsZero() {
					if nowMinusintervalOneDay.Before(created) && now.After(created) {
						result = append(result, v)
					}
				} else if startTime.Before(created) && EndTime.After(created) {
					result = append(result, v)
				}
			}
		}
		return result
	}

	url := client.getURL("subscriptions/content", nil)
	mux.HandleFunc(url.Path, func(w http.ResponseWriter, r *http.Request) {
		EnforceMethod(t, r, "GET")
		contentType := EnforceAndReturnContentType(t, r)

		startTimeStr := r.URL.Query().Get("startTime")
		endTimeStr := r.URL.Query().Get("endTime")
		startTime := time.Time{}
		endTime := time.Time{}
		if startTimeStr != "" && endTimeStr != "" {
			startTime = EnforceAndReturnTime(t, r, startTimeStr)
			endTime = EnforceAndReturnTime(t, r, endTimeStr)
		}

		response := filterStore(&store, contentType, startTime, endTime)

		maxCount := 3
		nextPage := r.URL.Query().Get("nextPage")
		if nextPage != "" {
			pageIndex, _ := strconv.Atoi(nextPage)
			lastPageIndex := len(response) - maxCount - (len(response) % maxCount)
			if 0 < pageIndex && pageIndex <= lastPageIndex && pageIndex%maxCount == 0 {
				response = response[pageIndex : pageIndex+maxCount]
				if pageIndex < lastPageIndex {
					nextPageURI, _ := url.Parse(r.URL.String())
					queryParams := nextPageURI.Query()
					queryParams.Set("nextPage", strconv.Itoa(pageIndex+maxCount))
					nextPageURI.RawQuery = queryParams.Encode()
					w.Header().Set("NextPageUri", nextPageURI.String())
				}
			}
		} else if len(response) > maxCount {
			response = response[0:maxCount]

			nextPageURI, _ := url.Parse(r.URL.String())
			queryParams := nextPageURI.Query()
			queryParams.Set("nextPage", strconv.Itoa(maxCount))
			nextPageURI.RawQuery = queryParams.Encode()
			w.Header().Set("NextPageUri", nextPageURI.String())
		}

		json.NewEncoder(w).Encode(response)
	})

	cases := []struct {
		ContentType ContentType
		StartTime   time.Time
		EndTime     time.Time
		Want        []Content
		WantError   error
	}{
		{
			ContentType: AuditAzureActiveDirectory,
			StartTime:   now.Add(-(intervalOneDay * 6)),
			EndTime:     now.Add(-(intervalOneDay * 5)),
			Want:        store[0:1],
			WantError:   nil,
		},
		{
			ContentType: AuditAzureActiveDirectory,
			StartTime:   now.Add(-(intervalOneDay * 4)),
			EndTime:     now.Add(-(intervalOneDay * 3)),
			Want:        []Content{},
			WantError:   nil,
		},
		{
			ContentType: AuditAzureActiveDirectory,
			StartTime:   now.Add(-(intervalOneDay * 6)),
			EndTime:     now.Add(-(intervalOneDay * 4)),
			Want:        nil,
			WantError:   ErrIntervalDay,
		},
		{
			ContentType: AuditAzureActiveDirectory,
			StartTime:   now.Add(-(intervalOneDay * 4)),
			EndTime:     now.Add(-(intervalOneDay * 6)),
			Want:        nil,
			WantError:   ErrIntervalNegative,
		},
		{
			ContentType: AuditAzureActiveDirectory,
			StartTime:   now.Add(-(intervalOneDay * 8)),
			EndTime:     now.Add(-(intervalOneDay * 7)),
			Want:        nil,
			WantError:   ErrIntervalWeek,
		},
		{
			ContentType: AuditAzureActiveDirectory,
			StartTime:   now.Add(-(intervalOneDay * 5)),
			EndTime:     time.Time{},
			Want:        nil,
			WantError:   ErrIntervalMismatch,
		},
		{
			ContentType: AuditAzureActiveDirectory,
			StartTime:   time.Time{},
			EndTime:     now.Add(-(intervalOneDay * 5)),
			Want:        nil,
			WantError:   ErrIntervalMismatch,
		},
		{
			ContentType: DLPAll,
			StartTime:   time.Time{},
			EndTime:     time.Time{},
			Want:        store[3:4],
			WantError:   nil,
		},
		{
			ContentType: AuditExchange,
			StartTime:   now.Add(-(intervalOneDay * 2)),
			EndTime:     now.Add(-(intervalOneDay * 1)),
			Want:        store[4:10],
			WantError:   nil,
		},
		{
			ContentType: AuditExchange,
			StartTime:   now.Add(-(intervalOneDay * 2)),
			EndTime:     now.Add(-(intervalOneDay * 1)),
			Want:        store[4:10],
			WantError:   nil,
		},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d.", idx+1), func(t *testing.T) {
			contents, err := client.Subscriptions.Content(context.Background(), &c.ContentType, c.StartTime, c.EndTime)
			testError(t, c.Want, c.WantError, err)
			if len(contents) == 0 && len(c.Want) == 0 {
				return
			}
			testDeep(t, contents, c.Want)
		})
	}
}

func TestAudit(t *testing.T) {

	client, mux, teardown := stubClient()
	defer teardown()

	store := map[string][]AuditRecord{
		"abc": {
			{ID: "qqqqqqq"},
		},
		"deg": {
			{ID: "123456"},
			{ID: "789012"},
		},
	}

	filterStore := func(c *map[string][]AuditRecord, contentID string) []AuditRecord {
		var result []AuditRecord
		for k, v := range *c {
			if k == contentID {
				result = append(result, v...)
			}
		}
		return result
	}

	url := client.getURL("audit/", nil)
	mux.HandleFunc(url.Path, func(w http.ResponseWriter, r *http.Request) {
		EnforceMethod(t, r, "GET")

		tokens := strings.Split(r.URL.Path, `/`)
		contentID := tokens[len(tokens)-1]
		response := filterStore(&store, contentID)

		json.NewEncoder(w).Encode(response)
	})

	cases := []struct {
		ContentID string
		Want      []AuditRecord
		WantError error
	}{
		{ContentID: "abc", Want: store["abc"], WantError: nil},
		{ContentID: "def", Want: []AuditRecord{}, WantError: nil},
		{ContentID: "deg", Want: store["deg"], WantError: nil},
		{ContentID: "", Want: nil, WantError: fmt.Errorf("ContentID must not be empty")},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d.", idx), func(t *testing.T) {
			records, err := client.Subscriptions.Audit(context.Background(), c.ContentID)
			testError(t, c.Want, c.WantError, err)
			if len(records) == 0 && len(c.Want) == 0 {
				return
			}
			testDeep(t, records, c.Want)
		})
	}
}
