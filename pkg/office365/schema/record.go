package schema

import (
	"encoding/json"
	"fmt"
)

// AuditRecord represents an event or action returned by Audit endpoint.
type AuditRecord struct {
	ID             *string             `json:"Id"`
	RecordType     *AuditLogRecordType `json:"RecordType"`
	CreationTime   *string             `json:"CreationTime"`
	Operation      *string             `json:"Operation"`
	OrganizationID *string             `json:"OrganizationId"`
	UserType       *UserType           `json:"UserType"`
	UserKey        *string             `json:"UserKey"`
	Workload       *string             `json:"Workload,omitempty"`
	ResultStatus   *string             `json:"ResultStatus,omitempty"`
	ObjectID       *string             `json:"ObjectId,omitempty"`
	UserID         *string             `json:"UserId"`
	ClientIP       *string             `json:"ClientIP"`
	Scope          *AuditLogScope      `json:"Scope,omitempty"`
}

// AuditLogRecordType identifies the type of AuditRecord.
// https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-schema#enum-auditlogrecordtype---type-edmint32
type AuditLogRecordType int

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

// GetRecordType returns the RecordType for the provided string.
func GetRecordType(s string) (*AuditLogRecordType, error) {
	literals := map[string]AuditLogRecordType{
		"ExchangeAdmin":                         ExchangeAdminType,
		"ExchangeItem":                          ExchangeItemType,
		"ExchangeItemGroup":                     ExchangeItemGroupType,
		"SharePoint":                            SharePointType,
		"SharePointFileOperation":               SharePointFileOperationType,
		"AzureActiveDirectory":                  AzureActiveDirectoryType,
		"AzureActiveDirectoryAccountLogon":      AzureActiveDirectoryAccountLogonType,
		"DataCenterSecurityCmdlet":              DataCenterSecurityCmdletType,
		"ComplianceDLPSharePoint":               ComplianceDLPSharePointType,
		"Sway":                                  SwayType,
		"ComplianceDLPExchange":                 ComplianceDLPExchangeType,
		"SharePointSharingOperation":            SharePointSharingOperationType,
		"AzureActiveDirectoryStsLogon":          AzureActiveDirectoryStsLogonType,
		"SecurityComplianceCenterEOPCmdlet":     SecurityComplianceCenterEOPCmdletType,
		"PowerBIAudit":                          PowerBIAuditType,
		"CRM":                                   CRMType,
		"Yammer":                                YammerType,
		"SkypeForBusinessCmdlets":               SkypeForBusinessCmdletsType,
		"Discovery":                             DiscoveryType,
		"MicrosoftTeams":                        MicrosoftTeamsType,
		"ThreatIntelligence":                    ThreatIntelligenceType,
		"MailSubmission":                        MailSubmissionType,
		"MicrosoftFlow":                         MicrosoftFlowType,
		"AeD":                                   AeDType,
		"MicrosoftStream":                       MicrosoftStreamType,
		"ComplianceDLPSharePointClassification": ComplianceDLPSharePointClassificationType,
		"Project":                               ProjectType,
		"SharePointListOperation":               SharePointListOperationType,
		"DataGovernance":                        DataGovernanceType,
		"SecurityComplianceAlerts":              SecurityComplianceAlertsType,
		"ThreatIntelligenceUrl":                 ThreatIntelligenceURLType,
		"SecurityComplianceInsights":            SecurityComplianceInsightsType,
		"WorkplaceAnalytics":                    WorkplaceAnalyticsType,
		"PowerAppsApp":                          PowerAppsAppType,
		"ThreatIntelligenceAtpContent":          ThreatIntelligenceAtpContentType,
		"TeamsHealthcare":                       TeamsHealthcareType,
		"DataInsightsRestApiAudit":              DataInsightsRestAPIAuditType,
		"SharePointListItemOperation":           SharePointListItemOperationType,
		"SharePointContentTypeOperation":        SharePointContentTypeOperationType,
		"SharePointFieldOperation":              SharePointFieldOperationType,
		"AirInvestigation":                      AirInvestigationType,
		"Quarantine":                            QuarantineType,
		"MicrosoftForms":                        MicrosoftFormsType,
	}
	t, ok := literals[s]
	if !ok {
		return nil, fmt.Errorf("record type invalid")
	}
	return &t, nil
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

// ContentType represents a type and source of aggregated actions and events
// generated by the Microsoft Office 365 Management Activity API.
type ContentType int

var contentTypeCount = 5

// ContentType enum.
const (
	AuditAzureActiveDirectory ContentType = iota
	AuditExchange
	AuditSharePoint
	AuditGeneral
	DLPAll
)

func (c ContentType) String() string {
	literals := map[ContentType]string{
		AuditAzureActiveDirectory: "Audit.AzureActiveDirectory",
		AuditExchange:             "Audit.Exchange",
		AuditSharePoint:           "Audit.SharePoint",
		AuditGeneral:              "Audit.General",
		DLPAll:                    "DLP.All",
	}
	return literals[c]
}

var contentTypes = map[string]ContentType{
	"Audit.AzureActiveDirectory": AuditAzureActiveDirectory,
	"Audit.Exchange":             AuditExchange,
	"Audit.SharePoint":           AuditSharePoint,
	"Audit.General":              AuditGeneral,
	"DLP.All":                    DLPAll,
}

// GetContentType returns the ContentType represented
// by the provided string literal.
func GetContentType(s string) (*ContentType, error) {
	if v, ok := contentTypes[s]; ok {
		return &v, nil
	}
	return nil, fmt.Errorf("ContentType invalid")
}

// GetContentTypes returns the list of ContentType.
func GetContentTypes() []ContentType {
	var result []ContentType
	for _, t := range contentTypes {
		result = append(result, t)
	}
	return result
}

// ContentTypeValid validates that a string is a valid ContentType.
func ContentTypeValid(s string) bool {
	if _, err := GetContentType(s); err != nil {
		return false
	}
	return true
}
