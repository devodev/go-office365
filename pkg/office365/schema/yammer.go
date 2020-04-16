package schema

// Yammer .
type Yammer struct {
	AuditRecord
	ActorUserID        *string `json:"ActorUserId,omitempty"`
	ActorYammerUserID  *int    `json:"ActorYammerUserId,omitempty"`
	DataExportType     *string `json:"DataExportType,omitempty"`
	FileID             *int    `json:"FileId,omitempty"`
	FileName           *string `json:"FileName,omitempty"`
	GroupName          *string `json:"GroupName,omitempty"`
	IsSoftDelete       *bool   `json:"IsSoftDelete,omitempty"`
	MessageID          *int    `json:"MessageId,omitempty"`
	YammerNetworkID    *int    `json:"YammerNetworkId,omitempty"`
	TargetUserID       *string `json:"TargetUserId,omitempty"`
	TargetYammerUserID *int    `json:"TargetYammerUserId,omitempty"`
	VersionID          *int    `json:"VersionId,omitempty"`
}
