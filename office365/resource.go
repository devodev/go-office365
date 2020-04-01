package office365

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

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
func (r *Resource) SetRequest(ct *ContentType, t time.Time) {
	r.Request = ResourceRequest{
		ContentType: ct,
		RequestTime: t,
	}
}

// SetResponse .
func (r *Resource) SetResponse(records []AuditRecord) {
	r.Response = ResourceResponse{records}
}

// ResourceRequest .
type ResourceRequest struct {
	ContentType *ContentType
	RequestTime time.Time
}

// ResourceResponse .
type ResourceResponse struct {
	Records []AuditRecord
}

// ResourceHandler is an interface for handling streamed resources.
type ResourceHandler interface {
	Handle(<-chan Resource)
}

// HumanReadableHandler implements the ResourceHandler interface.
// It prints a human readable formatted resource on the
// provided writer.
type HumanReadableHandler struct {
	writer io.Writer
}

// NewHumanReadableHandler returns a printer using the provided writer.
func NewHumanReadableHandler(w io.Writer) *HumanReadableHandler {
	return &HumanReadableHandler{w}
}

// Handle .
func (h HumanReadableHandler) Handle(in <-chan Resource) {
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

// JSONHandler implements the ResourceHandler interface.
// It writes json representation of a resource on the provided writer.
type JSONHandler struct {
	writer io.Writer
	logger *log.Logger
}

// NewJSONHandler returns a JSONHandler using the provided writer.
func NewJSONHandler(w io.Writer, l *log.Logger) *JSONHandler {
	return &JSONHandler{writer: w, logger: l}
}

// Handle .
func (h JSONHandler) Handle(in <-chan Resource) {
	for r := range in {
		for idx, e := range r.Errors {
			h.logger.Printf("[%s] Error%d: %s", r.Request.ContentType, idx, e.Error())
		}
		for _, a := range r.Response.Records {
			record := &JSONRecord{
				ContentType: r.Request.ContentType.String(),
				RequestTime: r.Request.RequestTime,
				Record:      a,
			}
			recordStr, err := json.Marshal(record)
			if err != nil {
				h.logger.Printf("error marshalling audit: %s\n", err)
				continue
			}
			fmt.Fprintln(h.writer, string(recordStr))
		}
	}
}

// JSONRecord is used for enriching AuditRecords with Request params.
type JSONRecord struct {
	ContentType string
	RequestTime time.Time
	Record      AuditRecord
}
