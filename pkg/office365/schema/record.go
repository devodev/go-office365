package schema

import (
	"encoding/json"
)

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

// AuditLogRecordType identifies the type of AuditRecord.
// https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-schema#enum-auditlogrecordtype---type-edmint32
type AuditLogRecordType int

// MarshalJSON marshals into a string.
func (t AuditLogRecordType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// AuditLogRecordType enum.
const (
	ExchangeAdminType AuditLogRecordType = iota + 1
	ExchangeItemType
	ExchangeItemGroupType
	SharePointType
	_
	SharePointFileOperationType
	_
	AzureActiveDirectoryType
	AzureActiveDirectoryAccountLogonType
	DataCenterSecurityCmdletType
	ComplianceDLPSharePointType
	SwayType
	ComplianceDLPExchangeType
	SharePointSharingOperationType
	AzureActiveDirectoryStsLogonType
	_
	_
	SecurityComplianceCenterEOPCmdletType
	_
	PowerBIAuditType
	CRMType
	YammerType
	SkypeForBusinessCmdletsType
	DiscoveryType
	MicrosoftTeamsType
	_
	_
	ThreatIntelligenceType
	MailSubmissionType
	MicrosoftFlowType
	AeDType
	MicrosoftStreamType
	ComplianceDLPSharePointClassificationType
	_
	ProjectType
	SharePointListOperationType
	_
	DataGovernanceType
	_
	SecurityComplianceAlertsType
	ThreatIntelligenceURLType
	SecurityComplianceInsightsType
	_
	WorkplaceAnalyticsType
	PowerAppsAppType
	_
	ThreatIntelligenceAtpContentType
	_
	TeamsHealthcareType
	_
	_
	DataInsightsRestAPIAuditType
	_
	SharePointListItemOperationType
	SharePointContentTypeOperationType
	SharePointFieldOperationType
	AirInvestigationType = iota + 8
	QuarantineType
	MicrosoftFormsType
)

func (t AuditLogRecordType) String() string {
	literals := map[AuditLogRecordType]string{
		ExchangeAdminType:                         "ExchangeAdmin",
		ExchangeItemType:                          "ExchangeItem",
		ExchangeItemGroupType:                     "ExchangeItemGroup",
		SharePointType:                            "SharePoint",
		SharePointFileOperationType:               "SharePointFileOperation",
		AzureActiveDirectoryType:                  "AzureActiveDirectory",
		AzureActiveDirectoryAccountLogonType:      "AzureActiveDirectoryAccountLogon",
		DataCenterSecurityCmdletType:              "DataCenterSecurityCmdlet",
		ComplianceDLPSharePointType:               "ComplianceDLPSharePoint",
		SwayType:                                  "Sway",
		ComplianceDLPExchangeType:                 "ComplianceDLPExchange",
		SharePointSharingOperationType:            "SharePointSharingOperation",
		AzureActiveDirectoryStsLogonType:          "AzureActiveDirectoryStsLogon",
		SecurityComplianceCenterEOPCmdletType:     "SecurityComplianceCenterEOPCmdlet",
		PowerBIAuditType:                          "PowerBIAudit",
		CRMType:                                   "CRM",
		YammerType:                                "Yammer",
		SkypeForBusinessCmdletsType:               "SkypeForBusinessCmdlets",
		DiscoveryType:                             "Discovery",
		MicrosoftTeamsType:                        "MicrosoftTeams",
		ThreatIntelligenceType:                    "ThreatIntelligence",
		MailSubmissionType:                        "MailSubmission",
		MicrosoftFlowType:                         "MicrosoftFlow",
		AeDType:                                   "AeD",
		MicrosoftStreamType:                       "MicrosoftStream",
		ComplianceDLPSharePointClassificationType: "ComplianceDLPSharePointClassification",
		ProjectType:                               "Project",
		SharePointListOperationType:               "SharePointListOperation",
		DataGovernanceType:                        "DataGovernance",
		SecurityComplianceAlertsType:              "SecurityComplianceAlerts",
		ThreatIntelligenceURLType:                 "ThreatIntelligenceUrl",
		SecurityComplianceInsightsType:            "SecurityComplianceInsights",
		WorkplaceAnalyticsType:                    "WorkplaceAnalytics",
		PowerAppsAppType:                          "PowerAppsApp",
		ThreatIntelligenceAtpContentType:          "ThreatIntelligenceAtpContent",
		TeamsHealthcareType:                       "TeamsHealthcare",
		DataInsightsRestAPIAuditType:              "DataInsightsRestApiAudit",
		SharePointListItemOperationType:           "SharePointListItemOperation",
		SharePointContentTypeOperationType:        "SharePointContentTypeOperation",
		SharePointFieldOperationType:              "SharePointFieldOperation",
		AirInvestigationType:                      "AirInvestigation",
		QuarantineType:                            "Quarantine",
		MicrosoftFormsType:                        "MicrosoftForms",
	}
	return literals[t]
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

func (t UserType) String() string {
	literals := map[UserType]string{
		Regular:          "Regular",
		Reserved:         "Reserved",
		Admin:            "Admin",
		DcAdmin:          "DcAdmin",
		System:           "System",
		Application:      "Application",
		ServicePrincipal: "ServicePrincipal",
		CustomPolicy:     "CustomPolicy",
		SystemPolicy:     "SystemPolicy",
	}
	return literals[t]
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

func (s AuditLogScope) String() string {
	literals := map[AuditLogScope]string{
		Online: "Online",
		Onprem: "Onprem",
	}
	return literals[s]
}
