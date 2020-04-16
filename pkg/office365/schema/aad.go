package schema

import "encoding/json"

// AzureActiveDirectoryBase .
type AzureActiveDirectoryBase struct {
	AzureActiveDirectoryEventType *AzureActiveDirectoryEventType `json:"AzureActiveDirectoryEventType"`
	ExtendedProperties            []NameValuePair                `json:"ExtendedProperties,omitempty"`
	ModifiedProperties            []string                       `json:"ModifiedProperties,omitempty"`
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

func (t AzureActiveDirectoryEventType) String() string {
	literals := map[AzureActiveDirectoryEventType]string{
		AccountLogon:               "AccountLogon",
		AzureApplicationAuditEvent: "AzureApplicationAuditEvent",
	}
	return literals[t]
}

// AzureActiveDirectoryAccountLogon .
type AzureActiveDirectoryAccountLogon struct {
	AuditRecord
	Application *string `json:"Application,omitempty"`
	Client      *string `json:"Client,omitempty"`
	LoginStatus *int    `json:"LoginStatus"`
	UserDomain  *string `json:"UserDomain"`
}

// AzureActiveDirectory .
type AzureActiveDirectory struct {
	Actor           []IdentityTypeValuePair `json:"Actor,omitempty"`
	ActorContextID  *string                 `json:"ActorContextId,omitempty"`
	ActorIPAddress  *string                 `json:"ActorIpAddress,omitempty"`
	InterSystemsID  *string                 `json:"InterSystemsId,omitempty"`
	IntraSystemsID  *string                 `json:"IntraSystemsId,omitempty"`
	SupportTicketID *string                 `json:"SupportTicketId,omitempty"`
	Target          []IdentityTypeValuePair `json:"Target,omitempty"`
	TargetContextID *string                 `json:"TargetContextId,omitempty"`
}

// IdentityTypeValuePair .
type IdentityTypeValuePair struct {
	ID   *string       `json:"ID"`
	Type *IdentityType `json:"Type"`
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

func (t IdentityType) String() string {
	literals := map[IdentityType]string{
		Claim: "Claim",
		Name:  "Name",
		Other: "Other",
		PUID:  "PUID",
		SPN:   "SPN",
		UPN:   "UPN",
	}
	return literals[t]
}

// AzureActiveDirectorySTSLogon .
type AzureActiveDirectorySTSLogon struct {
	AuditRecord
	ApplicationID *string `json:"ApplicationId,omitempty"`
	Client        *string `json:"Client,omitempty"`
	LogonError    *string `json:"LogonError,omitempty"`
}
