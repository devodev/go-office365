package schema

import "encoding/json"

// DataCenterSecurityBase .
type DataCenterSecurityBase struct {
	DataCenterSecurityEventType DataCenterSecurityEventType `json:"DataCenterSecurityEventType"`
}

// DataCenterSecurityEventType  .
type DataCenterSecurityEventType int

// MarshalJSON marshals into a string.
func (t DataCenterSecurityEventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// DataCenterSecurityEventType enum.
const (
	DataCenterSecurityCmdletAuditEvent DataCenterSecurityEventType = iota
)

func (t DataCenterSecurityEventType) String() string {
	literals := map[DataCenterSecurityEventType]string{
		DataCenterSecurityCmdletAuditEvent: "DataCenterSecurityCmdletAuditEvent",
	}
	return literals[t]
}

// DataCenterSecurityCmdlet .
type DataCenterSecurityCmdlet struct {
	AuditRecord
	StartTime             *string `json:"StartTime"`
	EffectiveOrganization *string `json:"EffectiveOrganization"`
	ElevationTime         *string `json:"ElevationTime"`
	ElevationApprover     *string `json:"ElevationApprover"`
	ElevationApprovedTime *string `json:"ElevationApprovedTime,omitempty"`
	ElevationRequestID    *string `json:"ElevationRequestId"`
	ElevationRole         *string `json:"ElevationRole,omitempty"`
	ElevationDuration     *int    `json:"ElevationDuration"`
	GenericInfo           *string `json:"GenericInfo,omitempty"`
}
