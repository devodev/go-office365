package office365

import "encoding/json"

// AuditRecord represents an event or action returned by Audit endpoint.
type AuditRecord struct {
	ID             string             `json:"Id"`
	RecordType     AuditLogRecordType `json:"RecordType"`
	CreationTime   string             `json:"CreationTime"`
	Operation      string             `json:"Operation"`
	OrganizationID string             `json:"OrganizationId"`
	UserType       UserType           `json:"UserType"`
	UserKey        string             `json:"UserKey"`
	Workload       string             `json:"Workload,omitempty"`
	ResultStatus   string             `json:"ResultStatus,omitempty"`
	ObjectID       string             `json:"ObjectId,omitempty"`
	UserID         string             `json:"UserId"`
	ClientIP       string             `json:"ClientIP"`
	Scope          AuditLogScope      `json:"Scope,omitempty"`
}

// MarshalJSON .
func (r AuditRecord) MarshalJSON() ([]byte, error) {
	return nil, nil
}

// AuditLogRecordType identifies the type of AuditRecord.
// https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-schema#enum-auditlogrecordtype---type-edmint32
type AuditLogRecordType int

// MarshalJSON marshals into a string.
func (t AuditLogRecordType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// AuditLogRecordType enum.
const (
	ExchangeAdmin AuditLogRecordType = iota + 1
	ExchangeItem
	ExchangeItemGroup
	SharePoint
	SharePointFileOperation
	AzureActiveDirectory
	AzureActiveDirectoryAccountLogon
	DataCenterSecurityCmdlet
	ComplianceDLPSharePoint
	Sway
	ComplianceDLPExchange
	SharePointSharingOperation
	AzureActiveDirectoryStsLogon
	SecurityComplianceCenterEOPCmdlet
	PowerBIAudit
	CRM
	Yammer
	SkypeForBusinessCmdlets
	Discovery
	MicrosoftTeams
	ThreatIntelligence
	MailSubmission
	MicrosoftFlow
	AeD
	MicrosoftStream
	ComplianceDLPSharePointClassification
	Project
	SharePointListOperation
	DataGovernance
	SecurityComplianceAlerts
	ThreatIntelligenceURL
	SecurityComplianceInsights
	WorkplaceAnalytics
	PowerAppsApp
	ThreatIntelligenceAtpContent
	TeamsHealthcare
	DataInsightsRestAPIAudit
	SharePointListItemOperation
	SharePointContentTypeOperation
	SharePointFieldOperation
	AirInvestigation
	Quarantine
	MicrosoftForms
)

var auditLogRecordTypeLiterals = []string{
	"ExchangeAdmin",
	"ExchangeItem",
	"ExchangeItemGroup",
	"SharePoint",
	"SharePointFileOperation",
	"AzureActiveDirectory",
	"AzureActiveDirectoryAccountLogon",
	"DataCenterSecurityCmdlet",
	"ComplianceDLPSharePoint",
	"Sway",
	"ComplianceDLPExchange",
	"SharePointSharingOperation",
	"AzureActiveDirectoryStsLogon",
	"SecurityComplianceCenterEOPCmdlet",
	"PowerBIAudit",
	"CRM",
	"Yammer",
	"SkypeForBusinessCmdlets",
	"Discovery",
	"MicrosoftTeams",
	"ThreatIntelligence",
	"MailSubmission",
	"MicrosoftFlow",
	"AeD",
	"MicrosoftStream",
	"ComplianceDLPSharePointClassification",
	"Project",
	"SharePointListOperation",
	"DataGovernance",
	"SecurityComplianceAlerts",
	"ThreatIntelligenceUrl",
	"SecurityComplianceInsights",
	"WorkplaceAnalytics",
	"PowerAppsApp",
	"ThreatIntelligenceAtpContent",
	"TeamsHealthcare",
	"DataInsightsRestApiAudit",
	"SharePointListItemOperation",
	"SharePointContentTypeOperation",
	"SharePointFieldOperation",
	"AirInvestigation",
	"Quarantine",
	"MicrosoftForms",
}

func (t AuditLogRecordType) String() string {
	return auditLogRecordTypeLiterals[t]
}

// UserType identifies the type of user in AuditRecord.
// https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-schema#enum-user-type---type-edmint32
type UserType int

// MarshalJSON marshals into a string.
func (t UserType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UserType enum.
const (
	Regular UserType = iota
	Reserved
	Admin
	DcAdmin
	System
	Application
	ServicePrincipal
	CustomPolicy
	SystemPolicy
)

var userTypeLiterals = []string{
	"Regular",
	"Reserved",
	"Admin",
	"DcAdmin",
	"System",
	"Application",
	"ServicePrincipal",
	"CustomPolicy",
	"SystemPolicy",
}

func (t UserType) String() string {
	return userTypeLiterals[t]
}

// AuditLogScope identifies the scope of an AuditRecord.
// https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-schema#auditlogscope
type AuditLogScope int

// MarshalJSON marshals into a string.
func (s AuditLogScope) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// AuditLogScope enum.
const (
	Online AuditLogScope = iota
	Onprem
)

var auditLogScopeLiterals = []string{
	"Online",
	"Onprem",
}

func (s AuditLogScope) String() string {
	return auditLogScopeLiterals[s]
}
