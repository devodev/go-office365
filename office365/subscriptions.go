package office365

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"
)

// RequestFormats are the time.Format vars we must follow when providing
// datetime params to an API endpoint.
var (
	RequestDateFormat          = "2006-01-02"
	RequestDatetimeFormat      = "2006-01-02T15:04"
	RequestDatetimeLargeFormat = "2006-01-02T15:04:05"
)

// error definition.
var (
	ErrContentTypeRequired = errors.New("ContentType queryParam is required")
	ErrIntervalMismatch    = errors.New("StartTime and EndTime must both be provided or not at all")
	ErrIntervalNegative    = errors.New("interval given is 0 or negative")
	ErrIntervalDay         = errors.New("interval given is more than 24 hours")
	ErrIntervalWeek        = errors.New("StartTime given is more than 7 days in the past")
)

// helpers.
var (
	intervalOneDay = time.Minute * 1440
)

// SubscriptionService .
type SubscriptionService service

// List returns the list of subscriptions.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#list-current-subscriptions
//
// List current subscriptions
// This operation returns a collection of the current subscriptions together with the associated webhooks.
func (s *SubscriptionService) List(ctx context.Context) ([]Subscription, error) {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)

	req, err := s.client.newRequest("GET", "subscriptions/list", params.Values, nil)
	if err != nil {
		return nil, err
	}

	var out []Subscription
	_, err = s.client.do(ctx, req, &out)
	return out, err
}

// Start will start a subscription for the specified content type.
// A payload can optionnaly be provided to enable a webhook
// that will send notifications periodically about available content.
// See below webhgook section for details.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#start-a-subscription
//
// This operation starts a subscription to the specified content type. If a subscription to the specified content type already exists, this operation is used to:
// - Update the properties of an active webhook.
// - Enable a webhook that was disabled because of excessive failed notifications.
// - Re-enable an expired webhook by specifying a later or null expiration date.
// - Remove a webhook.
//
// Webhook validation
//
// When the /start operation is called and a webhook is specified, we will send a validation notification
// to the specified webhook address to validate that an active listener can accept and process notifications.
//
// If we do not receive an HTTP 200 OK response, the subscription will not be created.
// Or, if /start is being called to add a webhook to an existing subscription and a response of HTTP 200 OK
// is not received, the webhook will not be added and the subscription will remain unchanged.
func (s *SubscriptionService) Start(ctx context.Context, ct *ContentType, webhook *Webhook) (*Subscription, error) {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)
	if err := params.AddContentType(ct); err != nil {
		return nil, err
	}

	var payload io.Reader
	if webhook != nil {
		data, err := json.Marshal(webhook)
		if err != nil {
			return nil, err
		}
		payload = bytes.NewBuffer(data)
	}

	req, err := s.client.newRequest("POST", "subscriptions/start", params.Values, payload)
	if err != nil {
		return nil, err
	}

	var out *Subscription
	_, err = s.client.do(ctx, req, &out)
	return out, err
}

// Stop stops a subscription for the provided ContentType.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#stop-a-subscription
//
// This operation stops a subscription to the specified content type.
// When a subscription is stopped, you will no longer receive notifications and you will not be able to retrieve available content.
// If the subscription is later restarted, you will have access to new content from that point forward.
// You will not be able to retrieve content that was available between the time the subscription was stopped and restarted.
func (s *SubscriptionService) Stop(ctx context.Context, ct *ContentType) error {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)
	if err := params.AddContentType(ct); err != nil {
		return err
	}

	req, err := s.client.newRequest("POST", "subscriptions/stop", params.Values, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// Content returns a list of content available for retrieval.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#list-available-content
//
// This operation lists the content currently available for retrieval for the specified content type.
// The content is an aggregation of actions and events harvested from multiple servers across multiple datacenters.
// The content will be listed in the order in which the aggregations become available, but the events and actions within
// the aggregations are not guaranteed to be sequential. An error is returned if the subscription status is disabled.
func (s *SubscriptionService) Content(ctx context.Context, ct *ContentType, startTime time.Time, endTime time.Time) ([]Content, error) {
	params := NewQueryParams()
	params.AddPubIdentifier(s.client.pubIdentifier)
	if err := params.AddContentType(ct); err != nil {
		return nil, err
	}
	if err := params.AddStartEndTime(startTime, endTime); err != nil {
		return nil, err
	}

	out := []Content{}
	var err error
	for {
		req, err := s.client.newRequest("GET", "subscriptions/content", params.Values, nil)
		if err != nil {
			return nil, err
		}

		var sub []Content
		response, err := s.client.do(ctx, req, &sub)
		if err != nil {
			return nil, err
		}
		out = append(out, sub...)

		nextPageURIStr := response.Header.Get("NextPageUri")
		if nextPageURIStr == "" {
			break
		}
		nextPageURI, err := url.Parse(nextPageURIStr)
		if err != nil {
			return nil, err
		}
		nextPage := nextPageURI.Query().Get("nextPage")
		if nextPage == "" {
			return nil, fmt.Errorf("nextPage is not present as queryParam of NextPageUri header")
		}
		params.Set("nextPage", nextPage)
	}
	return out, err
}

// Audit returns a list of events or actions.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#retrieving-content
// To retrieve a content blob, make a GET request against the corresponding content URI that is included
// in the list of available content and in the notifications sent to a webhook.
// The returned content will be a collection of one more actions or events in JSON format.
func (s *SubscriptionService) Audit(ctx context.Context, contentID string) ([]AuditRecord, error) {
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

// Watch .
func (s *SubscriptionService) Watch(ctx context.Context, fetcherCount int, intervalMinutes int) <-chan Resource {
	generatedChan := make(chan Resource)
	resultChan := make(chan Resource)

	for i := 0; i < fetcherCount; i++ {
		go s.fetcher(ctx, generatedChan, resultChan)
	}
	go s.resourceGenerator(ctx, intervalMinutes, generatedChan)

	return resultChan
}

func (s *SubscriptionService) resourceGenerator(ctx context.Context, intervalMinutes int, out chan Resource) {
	// TODO: change time.Second into time.Minute. This is to ease testing.
	tickerDur := time.Duration(intervalMinutes) * time.Second
	ticker := time.NewTicker(tickerDur)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(out)
			return
		case t := <-ticker.C:
			go func() {
				resource := Resource{}

				subscriptions, err := s.client.Subscriptions.List(ctx)
				if err != nil {
					resource.AddError(err)
					out <- resource
					return
				}

				// TODO: remove time.Minute
				startTime := t.Add(-(tickerDur + time.Minute))
				endTime := t

				for _, sub := range subscriptions {
					ct, err := GetContentType(sub.ContentType)
					if err != nil {
						fmt.Println(err)
						continue
					}
					resource.SetRequest(ct, startTime, endTime)
					out <- resource
				}
			}()
		}
	}
}

func (s *SubscriptionService) fetcher(ctx context.Context, in <-chan Resource, out chan Resource) {
	defer close(out)
	for r := range in {
		content, err := s.client.Subscriptions.Content(ctx, r.Request.ContentType, r.Request.StartTime, r.Request.EndTime)
		if err != nil {
			r.AddError(err)
			out <- r
			continue
		}

		var records []AuditRecord
		for _, c := range content {
			audits, err := s.client.Subscriptions.Audit(ctx, c.ContentID)
			if err != nil {
				r.AddError(err)
				out <- r
				continue
			}
			records = append(records, audits...)
		}
		r.SetResponse(records)
		out <- r
	}
}

// Resource .
type Resource struct {
	Request  ResourceRequest
	Response ResourceResponse
	Errors   []error
}

// AddError .
func (r *Resource) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

// SetRequest .
func (r *Resource) SetRequest(ct *ContentType, startTime time.Time, endTime time.Time) {
	r.Request = ResourceRequest{
		ContentType: ct,
		StartTime:   startTime,
		EndTime:     endTime,
	}
}

// SetResponse .
func (r *Resource) SetResponse(records []AuditRecord) {
	r.Response = ResourceResponse{records}
}

// ResourceRequest .
type ResourceRequest struct {
	ContentType *ContentType
	StartTime   time.Time
	EndTime     time.Time
}

// ResourceResponse .
type ResourceResponse struct {
	Records []AuditRecord
}

// QueryParams .
type QueryParams struct {
	url.Values
}

// NewQueryParams .
func NewQueryParams() *QueryParams {
	return &QueryParams{make(url.Values)}
}

// AddPubIdentifier .
func (p *QueryParams) AddPubIdentifier(pubIdentifier string) {
	if pubIdentifier != "" {
		p.Add("PublisherIdentifier", pubIdentifier)
	}
}

// AddContentType .
func (p *QueryParams) AddContentType(ct *ContentType) error {
	if &ct == nil {
		return ErrContentTypeRequired
	}
	p.Add("contentType", ct.String())
	return nil
}

// AddStartEndTime .
func (p *QueryParams) AddStartEndTime(startTime time.Time, endTime time.Time) error {
	oneOrMoreDatetime := !startTime.IsZero() || !endTime.IsZero()
	bothDatetime := !startTime.IsZero() && !endTime.IsZero()
	if oneOrMoreDatetime && !bothDatetime {
		return ErrIntervalMismatch
	}
	if bothDatetime {
		interval := endTime.Sub(startTime)
		if interval <= 0 {
			return ErrIntervalNegative
		}
		if interval > intervalOneDay {
			return ErrIntervalDay
		}
		if startTime.Before(time.Now().Add(-(intervalOneDay * 7))) {
			return ErrIntervalWeek
		}
		p.Add("startTime", startTime.Format(RequestDatetimeFormat))
		p.Add("endTime", endTime.Format(RequestDatetimeFormat))
	}
	return nil
}

// Subscription represents a response.
type Subscription struct {
	ContentType string  `json:"contentType"`
	Status      string  `json:"status"`
	Webhook     Webhook `json:"webhook"`
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
