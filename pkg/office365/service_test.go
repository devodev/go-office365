package office365

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func stubClient() (*Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := NewClient(nil, "test-tenandID", "", nil)
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
			t.Errorf("error occurred but different than WantError: %v != %v", gotError, wantError)
		}
		if wantError == nil {
			t.Errorf("error occurred but WantError is nil: %v", gotError)
		}
		return
	}
	if want == nil {
		t.Errorf("no error occurred but Want is nil")
		return
	}
}

func testDeep(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got\n%v\nbut want\n%v", got, want)
	}
}
