package schema

// MainInvestigation .
type MainInvestigation struct {
	InvestigationID   *string  `json:"InvestigationId,omitempty"`
	InvestigationName *string  `json:"InvestigationName,omitempty"`
	InvestigationType *string  `json:"InvestigationType,omitempty"`
	LastUpdateTimeUtc *string  `json:"LastUpdateTimeUtc,omitempty"`
	StartTimeUtc      *string  `json:"StartTimeUtc,omitempty"`
	Status            *string  `json:"Status,omitempty"`
	DeeplinkURL       *string  `json:"DeeplinkURL,omitempty"`
	Actions           []string `json:"Actions,omitempty"`
	Data              *string  `json:"Data,omitempty"`
}

// Actions .
type Actions struct {
	ID              *string  `json:"ID,omitempty"`
	ActionType      *string  `json:"ActionType,omitempty"`
	ActionStatus    *string  `json:"ActionStatus,omitempty"`
	ApprovedBy      *string  `json:"ApprovedBy,omitempty"`
	TimestampUtc    *string  `json:"TimestampUtc,omitempty"`
	ActionID        *string  `json:"ActionId,omitempty"`
	InvestigationID *string  `json:"InvestigationId,omitempty"`
	RelatedAlertIds []string `json:"RelatedAlertIds,omitempty"`
	StartTimeUtc    *string  `json:"StartTimeUtc,omitempty"`
	EndTimeUtc      *string  `json:"EndTimeUtc,omitempty"`
	Resource        *string  `json:"Resource,omitempty"`
	Entities        []string `json:"Entities,omitempty"`
	Related         *string  `json:"Related,omitempty"`
}

// EntityMailMessage .
type EntityMailMessage struct {
	Type              *string      `json:"Type,omitempty"`
	Files             []EntityFile `json:"Files,omitempty"`
	Recipient         *string      `json:"Recipient,omitempty"`
	Urls              []EntityURL  `json:"Urls,omitempty"`
	Sender            *string      `json:"Sender,omitempty"`
	SenderIP          *string      `json:"SenderIP,omitempty"`
	ReceivedDate      *string      `json:"ReceivedDate,omitempty"`
	NetworkMessageID  *string      `json:"NetworkMessageId,omitempty"`
	InternetMessageID *string      `json:"InternetMessageId,omitempty"`
	Subject           *string      `json:"Subject,omitempty"`
}

// EntityIP .
type EntityIP struct {
	Type    *string `json:"Type,omitempty"`
	Address *string `json:"Address,omitempty"`
}

// EntityURL .
type EntityURL struct {
	Type *string `json:"Type,omitempty"`
	URL  *string `json:"Url,omitempty"`
}

// EntityMailbox .
type EntityMailbox struct {
	Type                  *string `json:"Type,omitempty"`
	MailboxPrimaryAddress *string `json:"MailboxPrimaryAddress,omitempty"`
	DisplayName           *string `json:"DisplayName,omitempty"`
	Upn                   *string `json:"Upn,omitempty"`
}

// EntityFile .
type EntityFile struct {
	Type       *string  `json:"Type,omitempty"`
	Name       *string  `json:"Name,omitempty"`
	FileHashes []string `json:"FileHashes,omitempty"`
}

// EntityFileHash .
type EntityFileHash struct {
	Type      *string `json:"Type,omitempty"`
	Algorithm *string `json:"Algorithm,omitempty"`
	Value     *string `json:"Value,omitempty"`
}

// EntityMailCluster .
type EntityMailCluster struct {
	Type                  *string  `json:"Type,omitempty"`
	NetworkMessageIds     []string `json:"NetworkMessageIds,omitempty"`
	CountByDeliveryStatus []string `json:"CountByDeliveryStatus,omitempty"`
	CountByThreatType     []string `json:"CountByThreatType,omitempty"`
	Threats               []string `json:"Threats,omitempty"`
	Query                 *string  `json:"Query,omitempty"`
	QueryTime             *string  `json:"QueryTime,omitempty"`
	MailCount             *int     `json:"MailCount,omitempty"`
	Source                *string  `json:"Source,omitempty"`
}
