package office365

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

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
		nextPage := r.URL.Query().Get("nextpage")
		if nextPage != "" {
			pageIndex, _ := strconv.Atoi(nextPage)
			lastPageIndex := len(response) - maxCount - (len(response) % maxCount)
			if 0 < pageIndex && pageIndex <= lastPageIndex && pageIndex%maxCount == 0 {
				response = response[pageIndex : pageIndex+maxCount]
				if pageIndex < lastPageIndex {
					nextPageURI, _ := url.Parse(r.URL.String())
					queryParams := nextPageURI.Query()
					queryParams.Set("nextpage", strconv.Itoa(pageIndex+maxCount))
					nextPageURI.RawQuery = queryParams.Encode()
					w.Header().Set("NextPageUri", nextPageURI.String())
				}
			}
		} else if len(response) > maxCount {
			response = response[0:maxCount]

			nextPageURI, _ := url.Parse(r.URL.String())
			queryParams := nextPageURI.Query()
			queryParams.Set("nextpage", strconv.Itoa(maxCount))
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
			contents, err := client.Content.List(context.Background(), &c.ContentType, c.StartTime, c.EndTime)
			testError(t, c.Want, c.WantError, err)
			if len(contents) == 0 && len(c.Want) == 0 {
				return
			}
			testDeep(t, contents, c.Want)
		})
	}
}
