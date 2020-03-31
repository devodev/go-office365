package office365

import (
	"context"
	"fmt"
)

// AuditService .
type AuditService service

// List returns a list of events or actions.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#retrieving-content
// To retrieve a content blob, make a GET request against the corresponding content URI that is included
// in the list of available content and in the notifications sent to a webhook.
// The returned content will be a collection of one more actions or events in JSON format.
func (s *AuditService) List(ctx context.Context, contentID string) ([]AuditRecord, error) {
	if contentID == "" {
		return nil, fmt.Errorf("ContentID must not be empty")
	}
	path := fmt.Sprintf("audit/%s", contentID)
	req, err := s.client.newRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var out []AuditRecord
	_, err = s.client.do(ctx, req, &out)
	return out, err
}
