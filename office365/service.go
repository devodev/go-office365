package office365

import (
	"bytes"
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

	CreatedDatetimeFormat = "2006-01-02T15:04:05.999Z"
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

// service holds a pointer to the Client for service related
// methods to access Client methods, such as newRequest and do.
type service struct {
	client *Client
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

// ResourceHandler is an interface for handling streamed resources.
type ResourceHandler interface {
	Handle(<-chan Resource)
}

// Printer implements the ResourceHandler interface.
// It prints a human readable formatted resource on the
// provided writer.
type Printer struct {
	writer io.Writer
}

// NewPrinter returns a printer using the provided writer.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w}
}

// Handle .
func (h Printer) Handle(in <-chan Resource) {
	for r := range in {
		for idx, e := range r.Errors {
			fmt.Fprintf(h.writer, "[%s] Error%d: %s", r.Request.ContentType, idx, e.Error())
		}
		for _, a := range r.Response.Records {
			auditStr, err := json.Marshal(a)
			if err != nil {
				fmt.Fprintf(h.writer, "error marshalling audit: %s\n", err)
				continue
			}
			var out bytes.Buffer
			json.Indent(&out, auditStr, "", "\t")
			fmt.Fprintf(h.writer, "[%s]\n%s\n", r.Request.ContentType, out.String())
		}
	}
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
