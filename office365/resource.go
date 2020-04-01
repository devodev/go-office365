package office365

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
