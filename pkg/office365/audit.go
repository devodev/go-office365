package office365

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devodev/go-office365/v0/pkg/office365/schema"
)

// AuditService .
type AuditService service

// List returns a list of events or actions.
//
// Microsoft API Reference: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference#retrieving-content
// To retrieve a content blob, make a GET request against the corresponding content URI that is included
// in the list of available content and in the notifications sent to a webhook.
// The returned content will be a collection of one more actions or events in JSON format.
func (s *AuditService) List(ctx context.Context, contentID string) (*Response, []interface{}, error) {
	if contentID == "" {
		return nil, nil, fmt.Errorf("ContentID must not be empty")
	}
	path := fmt.Sprintf("audit/%s", contentID)
	req, err := s.client.newRequest("GET", path, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var records []json.RawMessage
	resp, err := s.client.do(ctx, req, &records)
	if err != nil {
		return resp, nil, err
	}

	var out []interface{}
	for _, raw := range records {
		var data interface{}

		var r schema.AuditRecord
		if err := json.Unmarshal(raw, &r); err != nil {
			return resp, nil, err
		}

		data = &r
		switch r.RecordType {
		case schema.ExchangeAdminType:
			var d schema.ExchangeAdmin
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.ExchangeItemType:
			var d schema.ExchangeItem
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.ExchangeItemGroupType:
		case schema.SharePointType:
			var d schema.Sharepoint
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.SharePointFileOperationType:
			var d schema.SharepointFileOperations
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.AzureActiveDirectoryType:
			var d schema.AzureActiveDirectory
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.AzureActiveDirectoryAccountLogonType:
			var d schema.AzureActiveDirectoryAccountLogon
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.DataCenterSecurityCmdletType:
			var d schema.DataCenterSecurityCmdlet
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.ComplianceDLPSharePointType:
		case schema.SwayType:
			var d schema.Sway
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.ComplianceDLPExchangeType:
		case schema.SharePointSharingOperationType:
			var d schema.SharepointSharing
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.AzureActiveDirectoryStsLogonType:
			var d schema.AzureActiveDirectorySTSLogon
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.SecurityComplianceCenterEOPCmdletType:
			var d schema.SecurityComplianceCenter
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.PowerBIAuditType:
			var d schema.PowerBI
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.CRMType:
		case schema.YammerType:
			var d schema.Yammer
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.SkypeForBusinessCmdletsType:
		case schema.DiscoveryType:
		case schema.MicrosoftTeamsType:
			var d schema.MicrosoftTeams
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.ThreatIntelligenceType:
			var d schema.ATP
			if err := json.Unmarshal(raw, &d); err == nil {
				data = &d
			}
		case schema.MailSubmissionType:
		case schema.MicrosoftFlowType:
		case schema.AeDType:
		case schema.MicrosoftStreamType:
		case schema.ComplianceDLPSharePointClassificationType:
		case schema.ProjectType:
		case schema.SharePointListOperationType:
		case schema.DataGovernanceType:
		case schema.SecurityComplianceAlertsType:
		case schema.ThreatIntelligenceURLType:
		case schema.SecurityComplianceInsightsType:
		case schema.WorkplaceAnalyticsType:
		case schema.PowerAppsAppType:
		case schema.ThreatIntelligenceAtpContentType:
		case schema.TeamsHealthcareType:
		case schema.DataInsightsRestAPIAuditType:
		case schema.SharePointListItemOperationType:
		case schema.SharePointContentTypeOperationType:
		case schema.SharePointFieldOperationType:
		case schema.AirInvestigationType:
		case schema.QuarantineType:
		case schema.MicrosoftFormsType:
		}
		out = append(out, data)
	}

	return resp, out, err
}
