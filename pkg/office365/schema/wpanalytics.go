package schema

// WorkplaceAnalytics .
type WorkplaceAnalytics struct {
	AuditRecord
	WpaUserRole        *string         `json:"WpaUserRole,omitempty"`
	ModifiedProperties []string        `json:"ModifiedProperties,omitempty"`
	OperationDetails   []NameValuePair `json:"OperationDetails,omitempty"`
}
