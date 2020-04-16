package office365

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

// Operations URL
// https://manage.office.com/api/v1.0/{tenant_id}/activity/feed/{operation}

// Useful Troubleshooting Doc
// https://github.com/MicrosoftDocs/office-365-management-api/blob/master/office-365-management-api/troubleshooting-the-office-365-management-activity-api.md

var (
	defaultBaseURL   = "https://manage.office.com"
	defaultVersion   = "v1.0"
	defaultUserAgent = "go-office365"
	defaultTimeout   = 5 * time.Second

	microsoftTokenURL = "https://login.windows.net/%s/oauth2/token?api-version=1.0"
)

var (
	// ErrBadRequest is a 400 http error.
	ErrBadRequest = errors.New("bad request")
	// ErrNotFound is a 404 http error.
	ErrNotFound = errors.New("not found")
)

// Credentials are used by OAuthClient.
type Credentials struct {
	ClientID     string
	ClientSecret string
	TenantDomain string
	TenantID     string
}

// OAuthClient returns an authenticated httpClient using the provided credentials.
func OAuthClient(ctx context.Context, c *Credentials) *http.Client {
	conf := &clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     fmt.Sprintf(microsoftTokenURL, c.TenantDomain),
		EndpointParams: url.Values{
			"resource": []string{defaultBaseURL},
		},
	}
	return conf.Client(ctx)
}

// A Client handles communication with the
// Microsoft Graph REST API.
type Client struct {
	BaseURL   *url.URL
	UserAgent string
	version   string

	client        *http.Client
	tenantID      string
	pubIdentifier string

	// inspired by go-github:
	// https://github.com/google/go-github/blob/d913de9ce1e8ed5550283b448b37b721b61cc3b3/github/github.go#L159
	// Reuse a single struct instead of allocating one for each service on the heap.
	common service

	Subscription *SubscriptionService
	Content      *ContentService
	Audit        *AuditService
}

// NewClient creates a Client using the provided httpClient.
// If nil is provided, a default httpClient with a default timeout value is created.
// Note that the default client has no way of authenticating itself against
// the Microsoft Office365 Management Activity  API.
// A convenience function is provided just for that: NewClientAuthenticated.
func NewClient(httpClient *http.Client, tenantID string, pubIdentifier string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	if pubIdentifier == "" {
		pubIdentifier = tenantID
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		BaseURL:       baseURL,
		UserAgent:     defaultUserAgent,
		version:       defaultVersion,
		client:        httpClient,
		tenantID:      tenantID,
		pubIdentifier: pubIdentifier,
	}
	c.common.client = c

	c.Subscription = (*SubscriptionService)(&c.common)
	c.Content = (*ContentService)(&c.common)
	c.Audit = (*AuditService)(&c.common)
	return c
}

// Version returns the client version.
func (c *Client) Version() string {
	return c.version
}

// NewClientAuthenticated returns an authenticated Client.
// pubIdentifier is used on Microsoft side to group queries
// together in terms of quotas and limitations.
func NewClientAuthenticated(c *Credentials, pubIdentifier string) *Client {
	oauthClient := OAuthClient(context.Background(), c)
	return NewClient(oauthClient, c.TenantID, pubIdentifier)
}

// newRequest generates a http.Request based on the method
// and endpoint provided. Default headers are also set here.
func (c *Client) newRequest(method, path string, params url.Values, payload io.Reader) (*http.Request, error) {
	url := c.getURL(path, params)
	req, err := http.NewRequest(method, url.String(), payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

// getURL returns a URL based on the client version and tenantID
// for the path provided.
func (c *Client) getURL(path string, params url.Values) *url.URL {
	return &url.URL{
		Scheme:   c.BaseURL.Scheme,
		Host:     c.BaseURL.Host,
		Path:     fmt.Sprintf("/api/%s/%s/activity/feed/%s", c.version, c.tenantID, path),
		RawQuery: params.Encode(),
	}
}

// do performs a roundtrip using the underlying client
// and returns an error, if any.
// It will also try to decode the body into the provided out interface.
// It returns the response and any error from decoding.
func (c *Client) do(ctx context.Context, req *http.Request, out interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{resp}

	err = CheckResponse(resp)

	if out != nil {
		decErr := json.NewDecoder(resp.Body).Decode(&out)
		if decErr == io.EOF {
			decErr = nil
		}
		err = decErr
	}
	return response, err
}

// CheckResponse validates the response returned from
// an API call and returns an error, if any.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, &errorResponse.Err)
	}
	return errorResponse
}

// Response encapsulates the http response received from
// a successful API call.
type Response struct {
	Response *http.Response
}

// ErrorResponse encapsulates the http response as well as the
// error returned in the body of an API call.
type ErrorResponse struct {
	Response *http.Response
	Err      *Error
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v. API Error: %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Response.Status, r.Err)
}

// Error represents the json object returned in the body
// of the response when an error is encountered.
type Error struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// ContentType represents a type and source of aggregated actions and events
// generated by the Microsoft Office 365 Management Activity API.
type ContentType int

var contentTypeCount = 5

// ContentType enum.
const (
	AuditAzureActiveDirectory ContentType = iota
	AuditExchange
	AuditSharePoint
	AuditGeneral
	DLPAll
)

func (c ContentType) String() string {
	literals := map[ContentType]string{
		AuditAzureActiveDirectory: "Audit.AzureActiveDirectory",
		AuditExchange:             "Audit.Exchange",
		AuditSharePoint:           "Audit.SharePoint",
		AuditGeneral:              "Audit.General",
		DLPAll:                    "DLP.All",
	}
	return literals[c]
}

// GetContentType returns the ContentType represented
// by the provided string literal.
func GetContentType(s string) (*ContentType, error) {
	for idx, v := range contentTypeLiterals {
		if v == s {
			ct := ContentType(idx)
			return &ct, nil
		}
	}
	return nil, fmt.Errorf("ContentType invalid")
}

// GetContentTypes returns the list of ContentType.
func GetContentTypes() []ContentType {
	var result []ContentType
	for idx := range contentTypeLiterals {
		ct := ContentType(idx)
		result = append(result, ct)
	}
	return result
}

// ContentTypeValid validates that a string is a valid ContentType.
func ContentTypeValid(s string) bool {
	if _, err := GetContentType(s); err != nil {
		return false
	}
	return true
}
