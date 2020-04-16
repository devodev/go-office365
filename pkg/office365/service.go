package office365

import (
	"errors"
	"net/url"
	"time"

	"github.com/devodev/go-office365/v0/pkg/office365/schema"
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
	intervalOneDay  = time.Minute * 1440
	intervalOneWeek = intervalOneDay * 7
)

// service holds a pointer to the Client for service related
// methods to access Client methods, such as newRequest and do.
type service struct {
	client *Client
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
func (p *QueryParams) AddContentType(ct *schema.ContentType) error {
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
