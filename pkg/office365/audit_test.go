package office365

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

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
			_, records, err := client.Audit.List(context.Background(), c.ContentID)
			testError(t, c.Want, c.WantError, err)
			if len(records) == 0 && len(c.Want) == 0 {
				return
			}
			testDeep(t, records, c.Want)
		})
	}
}