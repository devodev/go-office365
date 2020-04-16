package schema

// DLP .
type DLP struct {
	SharePointMetaData               *SharePointMetadata `json:"SharePointMetaData,omitempty"`
	ExchangeMetaData                 *ExchangeMetadata   `json:"ExchangeMetaData,omitempty"`
	ExceptionInfo                    *string             `json:"ExceptionInfo,omitempty"`
	PolicyDetails                    []PolicyDetails     `json:"PolicyDetails"`
	SensitiveInfoDetectionIsIncluded *bool               `json:"SensitiveInfoDetectionIsIncluded"`
}

// SharePointMetadata .
type SharePointMetadata struct {
	From                 *string `json:"From"`
	ItemCreationTime     *string `json:"itemCreationTime"`
	SiteCollectionGUID   *string `json:"SiteCollectionGuid"`
	SiteCollectionURL    *string `json:"SiteCollectionUrl"`
	FileName             *string `json:"FileName"`
	FileOwner            *string `json:"FileOwner"`
	FilePathURL          *string `json:"FilePathUrl"`
	DocumentLastModifier *string `json:"DocumentLastModifier"`
	DocumentSharer       *string `json:"DocumentSharer"`
	UniqueID             *string `json:"UniqueId"`
	LastModifiedTime     *string `json:"LastModifiedTime"`
}

// ExchangeMetadata .
type ExchangeMetadata struct {
	MessageID      *string  `json:"MessageID"`
	From           *string  `json:"From"`
	To             []string `json:"To,omitempty"`
	CC             []string `json:"CC,omitempty"`
	BCC            []string `json:"BCC,omitempty"`
	Subject        *string  `json:"Subject"`
	Sent           *string  `json:"Sent"`
	RecipientCount *int     `json:"RecipientCount"`
}

// PolicyDetails .
type PolicyDetails struct {
	PolicyID   *string `json:"PolicyId"`
	PolicyName *string `json:"PolicyName"`
	Rules      []Rules `json:"Rules"`
}

// Rules .
type Rules struct {
	RuleID            *string            `json:"RuleId"`
	RuleName          *string            `json:"RuleName"`
	Actions           []string           `json:"Actions,omitempty"`
	OverriddenActions []string           `json:"OverriddenActions,omitempty"`
	Severity          *string            `json:"Severity,omitempty"`
	RuleMode          *string            `json:"RuleMode"`
	ConditionsMatched *ConditionsMatched `json:"ConditionsMatched,omitempty"`
}

// ConditionsMatched .
type ConditionsMatched struct {
	SensitiveInformation []SensitiveInformation `json:"SensitiveInformation,omitempty"`
	DocumentProperties   []NameValuePair        `json:"DocumentProperties,omitempty"`
	OtherConditions      []NameValuePair        `json:"OtherConditions,omitempty"`
}

// SensitiveInformation .
type SensitiveInformation struct {
	Confidence                     *int                            `json:"Confidence"`
	Count                          *int                            `json:"Count"`
	SensitiveType                  *string                         `json:"SensitiveType"`
	SensitiveInformationDetections *SensitiveInformationDetections `json:"SensitiveInformationDetections,omitempty"`
}

// SensitiveInformationDetections .
type SensitiveInformationDetections struct {
	Detections       []Detections `json:"Detections"`
	ResultsTruncated *bool        `json:"ResultsTruncated"`
}

// Detections .
type Detections map[string]*string
