package schema

// Project .
type Project struct {
	AuditRecord
	Entity          *string `json:"Entity"`
	Action          *string `json:"Action"`
	OnBehalfOfResID *string `json:"OnBehalfOfResId,omitempty"`
}
