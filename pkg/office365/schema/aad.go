package schema

import "encoding/json"

// AzureActiveDirectoryBase .
type AzureActiveDirectoryBase struct {
	AzureActiveDirectoryEventType AzureActiveDirectoryEventType `json:"AzureActiveDirectoryEventType"`
	ExtendedProperties            []NameValuePair               `json:"ExtendedProperties,omitempty"`
	ModifiedProperties            []string                      `json:"ModifiedProperties,omitempty"`
}

// AzureActiveDirectoryEventType .
type AzureActiveDirectoryEventType int

// MarshalJSON marshals into a string.
func (t AzureActiveDirectoryEventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// AzureActiveDirectoryEventType enum.
const (
	AccountLogon AzureActiveDirectoryEventType = iota
	AzureApplicationAuditEvent
)

// AzureActiveDirectoryEventTypeLiterals .
var AzureActiveDirectoryEventTypeLiterals = []string{
	"AccountLogon",
	"AzureApplicationAuditEvent",
}

func (t AzureActiveDirectoryEventType) String() string {
	return AzureActiveDirectoryEventTypeLiterals[t]
}

// AzureActiveDirectoryAccountLogon .
type AzureActiveDirectoryAccountLogon struct {
	Application string `json:"Application,omitempty"`
	Client      string `json:"Client,omitempty"`
	LoginStatus int    `json:"LoginStatus"`
	UserDomain  string `json:"UserDomain"`
}
