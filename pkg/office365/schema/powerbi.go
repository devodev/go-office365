package schema

// PowerBI .
type PowerBI struct {
	AuditRecord
	AppName               *string                     `json:"AppName,omitempty"`
	DashboardName         *string                     `json:"DashboardName,omitempty"`
	DataClassification    *string                     `json:"DataClassification,omitempty"`
	DatasetName           *string                     `json:"DatasetName,omitempty"`
	MembershipInformation []MembershipInformationType `json:"MembershipInformation,omitempty"`
	OrgAppPermission      *string                     `json:"OrgAppPermission,omitempty"`
	ReportName            *string                     `json:"ReportName,omitempty"`
	SharingInformation    []SharingInformationType    `json:"SharingInformation,omitempty"`
	SwitchState           *string                     `json:"SwitchState,omitempty"`
	WorkSpaceName         *string                     `json:"WorkSpaceName,omitempty"`
}

// MembershipInformationType .
type MembershipInformationType struct {
	MemberEmail *string `json:"MemberEmail,omitempty"`
	Status      *string `json:"Status,omitempty"`
}

//SharingInformationType .
type SharingInformationType struct {
	RecipientEmail    *string `json:"RecipientEmail,omitempty"`
	RecipientName     *string `json:"RecipientName,omitempty"`
	ResharePermission *string `json:"ResharePermission,omitempty"`
}
