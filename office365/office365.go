package office365

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	Subscriptions *SubscriptionService
}

// service holds a pointer to the Client for service related
// methods to access Client methods, such as newRequest and do.
type service struct {
	client *Client
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

	c.Subscriptions = (*SubscriptionService)(&c.common)
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
func (c *Client) do(ctx context.Context, req *http.Request, out interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, ErrBadRequest
		case http.StatusNotFound:
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf(resp.Status)
		}
	}
	if out != nil {
		return resp, json.NewDecoder(resp.Body).Decode(&out)
	}
	return resp, nil
}
