package schema

// SecurityComplianceCenter .
type SecurityComplianceCenter struct {
	AuditRecord
	StartTime             *string `json:"StartTime,omitempty"`
	ClientRequestID       *string `json:"ClientRequestId,omitempty"`
	CmdletVersion         *string `json:"CmdletVersion,omitempty"`
	EffectiveOrganization *string `json:"EffectiveOrganization,omitempty"`
	UserServicePlan       *string `json:"UserServicePlan,omitempty"`
	ClientApplication     *string `json:"ClientApplication,omitempty"`
	Parameters            *string `json:"Parameters,omitempty"`
	NonPiiParameters      *string `json:"NonPiiParameters,omitempty"`
}

// SecurityComplianceAlerts .
type SecurityComplianceAlerts struct {
	AuditRecord
	AlertID       *string `json:"AlertId"`
	AlertType     *string `json:"AlertType"`
	Name          *string `json:"Name"`
	PolicyID      *string `json:"PolicyId,omitempty"`
	Status        *string `json:"Status,omitempty"`
	Severity      *string `json:"Severity,omitempty"`
	Category      *string `json:"Category,omitempty"`
	Source        *string `json:"Source,omitempty"`
	Comments      *string `json:"Comments,omitempty"`
	Data          *string `json:"Data,omitempty"`
	AlertEntityID *string `json:"AlertEntityId,omitempty"`
	EntityType    *string `json:"EntityType,omitempty"`
}
