// This is heavily inspired by go-github testing techniques. Thanks!
// https://github.com/google/go-github
package office365

import (
	"testing"
)

func TestClientDefaults(t *testing.T) {
	tenantID := ""
	pubIdentifier := ""

	client := NewClient(nil, tenantID, pubIdentifier)
	if client == nil {
		t.Fatal("something went terribly wrong")
	}

	if client.client.Timeout != defaultTimeout {
		t.Errorf(
			"timeout is not default value. got: %v want: %v",
			client.client.Timeout, defaultTimeout,
		)
	}
	baseURL := client.BaseURL.String()
	if baseURL != defaultBaseURL {
		t.Errorf("baseURL is not default value. got: %v want: %v", baseURL, defaultBaseURL)
	}
	version := client.Version()
	if version != defaultVersion {
		t.Errorf("Version is not default value. got: %v want: %v", version, defaultVersion)
	}
	if pubIdentifier == "" && client.pubIdentifier != client.tenantID {
		t.Errorf("pubIdentifier is not default value(tenantID). got: %v want: %v", client.pubIdentifier, client.tenantID)
	}
}
