package schema

import "encoding/json"

// Quarantine .
type Quarantine struct {
	RequestType      RequestType   `json:"RequestType,omitempty"`
	RequestSource    RequestSource `json:"RequestSource,omitempty"`
	NetworkMessageID string        `json:"NetworkMessageId,omitempty"`
	ReleaseTo        string        `json:"ReleaseTo,omitempty"`
}

// RequestType .
type RequestType int

// MarshalJSON marshals into a string.
func (t RequestType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// RequestType enum.
const (
	Preview RequestType = iota
	Delete
	Release
	Export
	ViewHeader
)

// RequestTypeLiterals .
var RequestTypeLiterals = []string{
	"Preview",
	"Delete",
	"Release",
	"Export",
	"ViewHeader",
}

func (t RequestType) String() string {
	return RequestTypeLiterals[t]
}

// RequestSource .
type RequestSource int

// MarshalJSON marshals into a string.
func (t RequestSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// RequestSource enum.
const (
	SCC RequestSource = iota
	Cmdlet
	URLlink
)

// RequestSourceLiterals .
var RequestSourceLiterals = []string{
	"SCC",
	"Cmdlet",
	"URLlink",
}

func (t RequestSource) String() string {
	return RequestSourceLiterals[t]
}
