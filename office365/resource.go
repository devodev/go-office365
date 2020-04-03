package office365

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

// ResourceHandler is an interface for handling streamed resources.
type ResourceHandler interface {
	Handle(<-chan ResourceAudits)
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
func (h HumanReadableHandler) Handle(in <-chan ResourceAudits) {
	for res := range in {
		auditStr, err := json.Marshal(res.AuditRecord)
		if err != nil {
			fmt.Fprintf(h.writer, "error marshalling audit: %s\n", err)
			continue
		}
		var out bytes.Buffer
		json.Indent(&out, auditStr, "", "\t")
		fmt.Fprintf(h.writer, "[%s]\n%s\n", res.ContentType.String(), out.String())
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
func (h JSONHandler) Handle(in <-chan ResourceAudits) {
	for res := range in {
		record := &JSONRecord{
			ContentType: res.ContentType.String(),
			RequestTime: res.RequestTime,
			Record:      res.AuditRecord,
		}
		recordStr, err := json.Marshal(record)
		if err != nil {
			h.logger.Printf("marshalling error: %s", err)
			continue
		}
		fmt.Fprintln(h.writer, string(recordStr))
	}
}

// JSONRecord is used for enriching AuditRecords with Request params.
type JSONRecord struct {
	ContentType string
	RequestTime time.Time
	Record      AuditRecord
}
