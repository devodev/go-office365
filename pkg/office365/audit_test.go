package office365

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/devodev/go-office365/v0/pkg/office365/schema"
)

func TestAudit(t *testing.T) {

	client, mux, teardown := stubClient()
	defer teardown()

	tp := schema.ComplianceDLPExchangeType
	store := map[string][]interface{}{
		"abc": {
			schema.AuditRecord{ID: String("qqqqqqq"), RecordType: &tp},
		},
		"deg": {
			schema.AuditRecord{ID: String("123456"), RecordType: &tp},
			schema.AuditRecord{ID: String("789012"), RecordType: &tp},
		},
	}

	filterStore := func(c *map[string][]interface{}, contentID string) []interface{} {
		var result []interface{}
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
		Want      []interface{}
		WantError error
	}{
		{ContentID: "abc", Want: store["abc"], WantError: nil},
		{ContentID: "def", Want: []interface{}{}, WantError: nil},
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
