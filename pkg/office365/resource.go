package office365

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/sirupsen/logrus"
)

// ResourceHandler is an interface for handling streamed resources.
type ResourceHandler interface {
	Handle(<-chan ResourceAudits, *logrus.Logger) error
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
func (h HumanReadableHandler) Handle(in <-chan ResourceAudits, l *logrus.Logger) error {
	for res := range in {
		auditStr, err := json.Marshal(res.AuditRecord)
		if err != nil {
			l.Errorf("marshalling audit: %s", err)
			continue
		}
		var out bytes.Buffer
		err = json.Indent(&out, auditStr, "", "\t")
		if err != nil {
			l.Errorf("indenting json audit: %s", err)
			continue
		}
		fmt.Fprintf(h.writer, "[%s]\n%s\n", res.ContentType.String(), out.String())
	}
	return nil
}

// JSONHandler implements the ResourceHandler interface.
// It writes json representation of a resource on the provided writer.
type JSONHandler struct {
	writer io.Writer
	logger *log.Logger
}

// NewJSONHandler returns a JSONHandler using the provided writer.
func NewJSONHandler(w io.Writer) *JSONHandler {
	return &JSONHandler{writer: w}
}

// Handle .
func (h JSONHandler) Handle(in <-chan ResourceAudits, l *logrus.Logger) error {
	for res := range in {
		record := &JSONRecord{
			ContentType: res.ContentType.String(),
			RequestTime: res.RequestTime,
			Record:      res.AuditRecord,
		}
		recordStr, err := json.Marshal(record)
		if err != nil {
			l.Errorf("marshalling: %s", err)
			continue
		}
		fmt.Fprintln(h.writer, string(recordStr))
	}
	return nil
}

// JSONRecord is used for enriching AuditRecords with Request params.
type JSONRecord struct {
	ContentType string
	RequestTime time.Time
	Record      AuditRecord
}
