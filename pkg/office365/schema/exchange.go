package schema

import (
	"encoding/json"
)

// ExchangeAdmin .
type ExchangeAdmin struct {
	AuditRecord
	ModifiedObjectResolvedName *string         `json:"ModifiedObjectResolvedName,omitempty"`
	Parameters                 []NameValuePair `json:"Parameters,omitempty"`
	ModifiedProperties         []string        `json:"ModifiedProperties,omitempty"`
	ExternalAccess             *bool           `json:"ExternalAccess"`
	OriginatingServer          *string         `json:"OriginatingServer,omitempty"`
	OrganizationName           *string         `json:"OrganizationName,omitempty"`
}

// ExchangeMailbox .
type ExchangeMailbox struct {
	LogonType                    *LogonType `json:"LogonType,omitempty"`
	InternalLogonType            *LogonType `json:"InternalLogonType,omitempty"`
	MailboxGUID                  *string    `json:"MailboxGuid,omitempty"`
	MailboxOwnerUPN              *string    `json:"MailboxOwnerUPN,omitempty"`
	MailboxOwnerSid              *string    `json:"MailboxOwnerSid,omitempty"`
	MailboxOwnerMasterAccountSid *string    `json:"MailboxOwnerMasterAccountSid,omitempty"`
	LogonUserSid                 *string    `json:"LogonUserSid,omitempty"`
	LogonUserDisplayName         *string    `json:"LogonUserDisplayName,omitempty"`
	ExternalAccess               *bool      `json:"ExternalAccess"`
	OriginatingServer            *string    `json:"OriginatingServer,omitempty"`
	OrganizationName             *string    `json:"OrganizationName,omitempty"`
	ClientInfoString             *string    `json:"ClientInfoString,omitempty"`
	ClientIPAddress              *string    `json:"ClientIPAddress,omitempty"`
	ClientMachineName            *string    `json:"ClientMachineName,omitempty"`
	ClientProcessName            *string    `json:"ClientProcessName,omitempty"`
	ClientVersion                *string    `json:"ClientVersion,omitempty"`
}

// ExchangeMailboxAuditGroupRecord .
type ExchangeMailboxAuditGroupRecord struct {
	Folder                           *ExchangeFolder  `json:"Folder,omitempty"`
	CrossMailboxOperations           *bool            `json:"CrossMailboxOperations,omitempty"`
	DestMailboxID                    *string          `json:"DestMailboxId,omitempty"`
	DestMailboxOwnerUPN              *string          `json:"DestMailboxOwnerUPN,omitempty"`
	DestMailboxOwnerSid              *string          `json:"DestMailboxOwnerSid,omitempty"`
	DestMailboxOwnerMasterAccountSid *string          `json:"DestMailboxOwnerMasterAccountSid,omitempty"`
	DestFolder                       *ExchangeFolder  `json:"DestFolder,omitempty"`
	Folders                          []ExchangeFolder `json:"Folders,omitempty"`
	AffectedItems                    []ExchangeItem   `json:"AffectedItems,omitempty"`
}

// ExchangeMailboxAuditRecord .
type ExchangeMailboxAuditRecord struct {
	Item                          *ExchangeItem `json:"Item,omitempty"`
	ModifiedProperties            []string      `json:"ModifiedProperties,omitempty"`
	SendAsUserSMTP                *string       `json:"SendAsUserSmtp,omitempty"`
	SendAsUserMailboxGUID         *string       `json:"SendAsUserMailboxGuid,omitempty"`
	SendOnBehalfOfUserSMTP        *string       `json:"SendOnBehalfOfUserSmtp,omitempty"`
	SendOnBehalfOfUserMailboxGUID *string       `json:"SendOnBehalfOfUserMailboxGuid,omitempty"`
}

// ExchangeItem .
type ExchangeItem struct {
	AuditRecord
	ID           *string         `json:"Id"`
	Subject      *string         `json:"Subject,omitempty"`
	ParentFolder *ExchangeFolder `json:"ParentFolder,omitempty"`
	Attachments  *string         `json:"Attachments,omitempty"`
}

// ExchangeFolder .
type ExchangeFolder struct {
	ID   *string `json:"Id"`
	Path *string `json:"Path,omitempty"`
}

// LogonType .
type LogonType int

// MarshalJSON marshals into a string.
func (t LogonType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// LogonType enum.
const (
	OwnerLT LogonType = iota
	AdminLT
	DelegatedLT
	TransportLT
	SystemServiceLT
	BestAccessLT
	DelegatedAdminLT
)

func (t LogonType) String() string {
	literals := map[LogonType]string{
		OwnerLT:          "Owner",
		AdminLT:          "Admin",
		DelegatedLT:      "Delegated",
		TransportLT:      "Transport",
		SystemServiceLT:  "SystemService",
		BestAccessLT:     "BestAccess",
		DelegatedAdminLT: "DelegatedAdmin",
	}
	return literals[t]
}
