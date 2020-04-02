package office365

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

	logger *GoLogger

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
func NewClient(httpClient *http.Client, tenantID string, pubIdentifier string, l *log.Logger) *Client {
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
		logger:        NewLogger(l),
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
func NewClientAuthenticated(c *Credentials, pubIdentifier string, l *log.Logger) *Client {
	oauthClient := OAuthClient(context.Background(), c)
	return NewClient(oauthClient, c.TenantID, pubIdentifier, l)
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

// Subscription represents a response.
type Subscription struct {
	ContentType string   `json:"contentType"`
	Status      string   `json:"status"`
	Webhook     *Webhook `json:"webhook"`
}

// Webhook represents both a response and a request payload.
type Webhook struct {
	Status     string `json:"status,omitempty"`
	Address    string `json:"address"`
	AuthID     string `json:"authId,omitempty"`
	Expiration string `json:"expiration,omitempty"`
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

var contentTypeLiterals = []string{
	"Audit.AzureActiveDirectory",
	"Audit.Exchange",
	"Audit.SharePoint",
	"Audit.General",
	"DLP.All",
}

func (c ContentType) String() string {
	return contentTypeLiterals[c]
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

// Content represents metadata needed for retreiving aggregated data.
type Content struct {
	ContentType       string `json:"contentType"`
	ContentID         string `json:"contentId"`
	ContentURI        string `json:"contentUri"`
	ContentCreated    string `json:"contentCreated"`
	ContentExpiration string `json:"contentExpiration"`
}

// AuditRecord represents an event or action returned by Audit endpoint.
type AuditRecord struct {
	ID             string             `json:"Id"`
	RecordType     AuditLogRecordType `json:"RecordType"`
	CreationTime   string             `json:"CreationTime"`
	Operation      string             `json:"Operation"`
	OrganizationID string             `json:"OrganizationId"`
	UserType       UserType           `json:"UserType"`
	UserKey        string             `json:"UserKey"`
	Workload       string             `json:"Workload,omitempty"`
	ResultStatus   string             `json:"ResultStatus,omitempty"`
	ObjectID       string             `json:"ObjectId,omitempty"`
	UserID         string             `json:"UserId"`
	ClientIP       string             `json:"ClientIP"`
	Scope          AuditLogScope      `json:"Scope,omitempty"`
}

// AuditLogRecordType identifies the type of AuditRecord.
// https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-schema#enum-auditlogrecordtype---type-edmint32
type AuditLogRecordType int

// AuditLogRecordType enum.
const (
	ExchangeAdmin AuditLogRecordType = iota + 1
	ExchangeItem
	ExchangeItemGroup
	SharePoint
	SharePointFileOperation
	AzureActiveDirectory
	AzureActiveDirectoryAccountLogon
	DataCenterSecurityCmdlet
	ComplianceDLPSharePoint
	Sway
	ComplianceDLPExchange
	SharePointSharingOperation
	AzureActiveDirectoryStsLogon
	SecurityComplianceCenterEOPCmdlet
	PowerBIAudit
	CRM
	Yammer
	SkypeForBusinessCmdlets
	Discovery
	MicrosoftTeams
	ThreatIntelligence
	MailSubmission
	MicrosoftFlow
	AeD
	MicrosoftStream
	ComplianceDLPSharePointClassification
	Project
	SharePointListOperation
	DataGovernance
	SecurityComplianceAlerts
	ThreatIntelligenceURL
	SecurityComplianceInsights
	WorkplaceAnalytics
	PowerAppsApp
	ThreatIntelligenceAtpContent
	TeamsHealthcare
	DataInsightsRestAPIAudit
	SharePointListItemOperation
	SharePointContentTypeOperation
	SharePointFieldOperation
	AirInvestigation
	Quarantine
	MicrosoftForms
)

var auditLogRecordTypeLiterals = []string{
	"ExchangeAdmin",
	"ExchangeItem",
	"ExchangeItemGroup",
	"SharePoint",
	"SharePointFileOperation",
	"AzureActiveDirectory",
	"AzureActiveDirectoryAccountLogon",
	"DataCenterSecurityCmdlet",
	"ComplianceDLPSharePoint",
	"Sway",
	"ComplianceDLPExchange",
	"SharePointSharingOperation",
	"AzureActiveDirectoryStsLogon",
	"SecurityComplianceCenterEOPCmdlet",
	"PowerBIAudit",
	"CRM",
	"Yammer",
	"SkypeForBusinessCmdlets",
	"Discovery",
	"MicrosoftTeams",
	"ThreatIntelligence",
	"MailSubmission",
	"MicrosoftFlow",
	"AeD",
	"MicrosoftStream",
	"ComplianceDLPSharePointClassification",
	"Project",
	"SharePointListOperation",
	"DataGovernance",
	"SecurityComplianceAlerts",
	"ThreatIntelligenceUrl",
	"SecurityComplianceInsights",
	"WorkplaceAnalytics",
	"PowerAppsApp",
	"ThreatIntelligenceAtpContent",
	"TeamsHealthcare",
	"DataInsightsRestApiAudit",
	"SharePointListItemOperation",
	"SharePointContentTypeOperation",
	"SharePointFieldOperation",
	"AirInvestigation",
	"Quarantine",
	"MicrosoftForms",
}

func (a AuditLogRecordType) String() string {
	return auditLogRecordTypeLiterals[a]
}

// UserType identifies the type of user in AuditRecord.
// https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-schema#enum-user-type---type-edmint32
type UserType int

// UserType enum.
const (
	Regular UserType = iota
	Reserved
	Admin
	DcAdmin
	System
	Application
	ServicePrincipal
	CustomPolicy
	SystemPolicy
)

var userTypeLiterals = []string{
	"Regular",
	"Reserved",
	"Admin",
	"DcAdmin",
	"System",
	"Application",
	"ServicePrincipal",
	"CustomPolicy",
	"SystemPolicy",
}

func (u UserType) String() string {
	return userTypeLiterals[u]
}

// AuditLogScope identifies the scope of an AuditRecord.
// https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-schema#auditlogscope
type AuditLogScope int

// AuditLogScope enum.
const (
	Online AuditLogScope = iota
	Onprem
)

var auditLogScopeLiterals = []string{
	"Online",
	"Onprem",
}

func (a AuditLogScope) String() string {
	return auditLogScopeLiterals[a]
}
