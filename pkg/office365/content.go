package office365

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/devodev/go-office365/v0/pkg/office365/schema"
)

// ContentService .
type ContentService service

// List returns a list of content available for retrieval.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#list-available-content
//
// This operation lists the content currently available for retrieval for the specified content type.
// The content is an aggregation of actions and events harvested from multiple servers across multiple datacenters.
// The content will be listed in the order in which the aggregations become available, but the events and actions within
// the aggregations are not guaranteed to be sequential. An error is returned if the subscription status is disabled.
func (s *ContentService) List(ctx context.Context, ct *schema.ContentType, startTime time.Time, endTime time.Time) ([]*Response, []Content, error) {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)
	if err := params.AddContentType(ct); err != nil {
		return nil, nil, err
	}
	if err := params.AddStartEndTime(startTime, endTime); err != nil {
		return nil, nil, err
	}

	out := []Content{}
	responses := []*Response{}
	for {
		req, err := s.client.newRequest("GET", "subscriptions/content", params.Values, nil)
		if err != nil {
			return responses, nil, err
		}

		var sub []Content
		response, err := s.client.do(ctx, req, &sub)
		if err != nil {
			return responses, nil, err
		}
		responses = append(responses, response)
		out = append(out, sub...)

		nextPageURIStr := response.Response.Header.Get("NextPageUri")
		if nextPageURIStr == "" {
			break
		}
		nextPageURI, err := url.ParseRequestURI(nextPageURIStr)
		if err != nil {
			return responses, nil, err
		}
		nextPage := nextPageURI.Query().Get("nextpage")
		if nextPage == "" {
			return responses, nil, fmt.Errorf("nextpage is not present as queryParam of NextPageUri header")
		}
		params.Set("nextpage", nextPage)
	}
	return responses, out, nil
}

// Content represents metadata needed for retreiving aggregated data.
type Content struct {
	ContentType       string `json:"contentType"`
	ContentID         string `json:"contentId"`
	ContentURI        string `json:"contentUri"`
	ContentCreated    string `json:"contentCreated"`
	ContentExpiration string `json:"contentExpiration"`
}
