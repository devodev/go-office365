package schema

import "encoding/json"

// MicrosoftTeams .
type MicrosoftTeams struct {
	AuditRecord
	MessageID       *string                `json:"MessageId,omitempty"`
	Members         []MicrosoftTeamsMember `json:"Members,omitempty"`
	TeamName        *string                `json:"TeamName,omitempty"`
	TeamGUID        *string                `json:"TeamGuid,omitempty"`
	ChannelType     *string                `json:"ChannelType,omitempty"`
	ChannelName     *string                `json:"ChannelName,omitempty"`
	ChannelGUID     *string                `json:"ChannelGuid,omitempty"`
	ExtraProperties []KeyValuePair         `json:"ExtraProperties,omitempty"`
	AddOnType       *AddOnType             `json:"AddOnType,omitempty"`
	AddonName       *string                `json:"AddonName,omitempty"`
	AddOnGUID       *string                `json:"AddOnGuid,omitempty"`
	TabType         *string                `json:"TabType,omitempty"`
	Name            *string                `json:"Name,omitempty"`
	OldValue        *string                `json:"OldValue,omitempty"`
	NewValue        *string                `json:"NewValue,omitempty"`
}

// MicrosoftTeamsMember .
type MicrosoftTeamsMember struct {
	UPN         *string         `json:"UPN,omitempty"`
	Role        *MemberRoleType `json:"Role,omitempty"`
	DisplayName *string         `json:"DisplayName,omitempty"`
}

// MemberRoleType  .
type MemberRoleType int

// MarshalJSON marshals into a string.
func (t MemberRoleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// MemberRoleType enum.
const (
	MemberRT MemberRoleType = iota
	OwnerRT
	GuestRT
)

func (t MemberRoleType) String() string {
	literals := map[MemberRoleType]string{
		MemberRT: "Member",
		OwnerRT:  "Owner",
		GuestRT:  "Guest",
	}
	return literals[t]
}

// AddOnType  .
type AddOnType int

// MarshalJSON marshals into a string.
func (t AddOnType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// AddOnType enum.
const (
	Bot AddOnType = iota
	Connector
	Tab
)

func (t AddOnType) String() string {
	literals := map[AddOnType]string{
		Bot:       "Bot",
		Connector: "Connector",
		Tab:       "Tab",
	}
	return literals[t]
}

// KeyValuePair .
type KeyValuePair struct {
	Key   *string `json:"Key,omitempty"`
	Value *string `json:"Value,omitempty"`
}
