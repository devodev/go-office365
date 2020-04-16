package office365

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

// ResourceHandler is an interface for handling streamed resources.
type ResourceHandler interface {
	Handle(<-chan ResourceAudits) error
}

// JSONHandler implements the ResourceHandler interface.
// It writes json representation of a resource on the provided writer.
type JSONHandler struct {
	writer io.Writer
	logger *logrus.Logger
	indent bool
}

// NewJSONHandler returns a JSONHandler using the provided writer.
func NewJSONHandler(w io.Writer, l *logrus.Logger, indent bool) *JSONHandler {
	return &JSONHandler{w, l, indent}
}

// Handle .
func (h JSONHandler) Handle(in <-chan ResourceAudits) error {
	for res := range in {
		record := &JSONRecord{
			ContentType: res.ContentType.String(),
			RequestTime: res.RequestTime,
			Record:      res.AuditRecord,
		}
		recordStr, err := json.Marshal(record)
		if err != nil {
			h.logger.Error(err)
			continue
		}
		if !h.indent {
			fmt.Fprintln(h.writer, string(recordStr))
			continue
		}
		var out bytes.Buffer
		err = json.Indent(&out, recordStr, "", "\t")
		if err != nil {
			h.logger.Error(err)
			continue
		}
		fmt.Fprintln(h.writer, out.String())
	}
	return nil
}

// JSONRecord is used for enriching AuditRecords with Request params.
type JSONRecord struct {
	ContentType string
	RequestTime time.Time
	Record      interface{}
}
