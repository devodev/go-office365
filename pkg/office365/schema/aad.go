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

// AzureActiveDirectory .
type AzureActiveDirectory struct {
	Actor           []IdentityTypeValuePair `json:"Actor,omitempty"`
	ActorContextID  string                  `json:"ActorContextId,omitempty"`
	ActorIPAddress  string                  `json:"ActorIpAddress,omitempty"`
	InterSystemsID  string                  `json:"InterSystemsId,omitempty"`
	IntraSystemsID  string                  `json:"IntraSystemsId,omitempty"`
	SupportTicketID string                  `json:"SupportTicketId,omitempty"`
	Target          []IdentityTypeValuePair `json:"Target,omitempty"`
	TargetContextID string                  `json:"TargetContextId,omitempty"`
}

// IdentityTypeValuePair .
type IdentityTypeValuePair struct {
	ID   string       `json:"ID"`
	Type IdentityType `json:"Type"`
}

// IdentityType .
type IdentityType int

// MarshalJSON marshals into a string.
func (t IdentityType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// IdentityType enum.
const (
	Claim IdentityType = iota
	Name
	Other
	PUID
	SPN
	UPN
)

// IdentityTypeLiterals .
var IdentityTypeLiterals = []string{
	"Claim",
	"Name",
	"Other",
	"PUID",
	"SPN",
	"UPN",
}

func (t IdentityType) String() string {
	return IdentityTypeLiterals[t]
}

// AzureActiveDirectorySTSLogon .
type AzureActiveDirectorySTSLogon struct {
	ApplicationID string `json:"ApplicationId,omitempty"`
	Client        string `json:"Client,omitempty"`
	LogonError    string `json:"LogonError,omitempty"`
}
