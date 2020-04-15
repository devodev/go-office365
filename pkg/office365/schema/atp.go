package schema

import "encoding/json"

// ATP .
type ATP struct {
	AttachmentData    []AttachmentData `json:"AttachmentData,omitempty"`
	DetectionType     string           `json:"DetectionType"`
	DetectionMethod   string           `json:"DetectionMethod"`
	InternetMessageID string           `json:"InternetMessageId"`
	NetworkMessageID  string           `json:"NetworkMessageId"`
	P1Sender          string           `json:"P1Sender"`
	P2Sender          string           `json:"P2Sender"`
	Policy            Policy           `json:"Policy"`
	PolicyAction      PolicyAction     `json:"PolicyAction"`
	Recipients        []string         `json:"Recipients"`
	SenderIP          string           `json:"SenderIp"`
	Subject           string           `json:"Subject"`
	Verdict           string           `json:"Verdict"`
	MessageTime       string           `json:"MessageTime"`
	EventDeepLink     string           `json:"EventDeepLink"`
}

// AttachmentData .
type AttachmentData struct {
	FileName      string      `json:"FileName"`
	FileType      string      `json:"FileType"`
	FileVerdict   FileVerdict `json:"FileVerdict"`
	MalwareFamily string      `json:"MalwareFamily,omitempty"`
	SHA256        string      `json:"SHA256"`
}

// FileVerdict .
type FileVerdict int

// MarshalJSON marshals into a string.
func (t FileVerdict) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// FileVerdict enum.
const (
	Pending FileVerdict = iota - 3
	Timeout
	Error
	Good
	Bad
)

// FileVerdictLiterals .
var FileVerdictLiterals = []string{
	"Pending",
	"Timeout",
	"Error",
	"Good",
	"Bad",
}

func (t FileVerdict) String() string {
	return FileVerdictLiterals[t]
}

// Policy .
type Policy int

// MarshalJSON marshals into a string.
func (t Policy) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// Policy enum.
const (
	AntiSpamHSPM Policy = iota + 1
	AntiSpamSPM
	AntiSpamBulk
	AntiSpamPHSH
	AntiPhishDIMP
	AntiPhishUIMP
	AntiPhishSPOOF
	AntiPhishGIMP
	AntiMalwareAMP
	SafeAttachmentSAP
	ExchangeTransport
	AntiMalwareZAPM
	AntiPhishZAPP
	AntiPhishZAPS
)

// PolicyLiterals .
var PolicyLiterals = []string{
	"Anti-spam, HSPM",
	"Anti-spam, SPM",
	"Anti-spam, Bulk",
	"Anti-spam, PHSH",
	"Anti-phish, DIMP",
	"Anti-phish, UIMP",
	"Anti-phish, SPOOF",
	"Anti-phish, GIMP",
	"Anti-malware, AMP",
	"Safe attachment, SAP",
	"Exchange transport",
	"Anti-malware, ZAPM",
	"Anti-phish, ZAPP",
	"Anti-phish, ZAPS",
}

func (t Policy) String() string {
	return PolicyLiterals[t]
}

// PolicyAction .
type PolicyAction int

// MarshalJSON marshals into a string.
func (t PolicyAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// PolicyAction enum.
const (
	MoveToJMFPA PolicyAction = iota
	AddXHeaderPA
	ModifySubjectPA
	RedirectPA
	DeletePA
	QuarantinePA
	NoActionPA
	BccMessagePA
	ReplaceAttachmentPA
)

// PolicyActionLiterals .
var PolicyActionLiterals = []string{
	"MoveToJMF",
	"AddXHeader",
	"ModifySubject",
	"Redirect",
	"Delete",
	"Quarantine",
	"NoAction",
	"BccMessage",
	"ReplaceAttachment",
}

func (t PolicyAction) String() string {
	return PolicyActionLiterals[t]
}

// URLTimeOfClickEvents .
type URLTimeOfClickEvents struct {
	UserID         string         `json:"UserId"`
	AppName        string         `json:"AppName"`
	URLClickAction URLClickAction `json:"URLClickAction"`
	SourceID       string         `json:"SourceId"`
	TimeOfClick    string         `json:"TimeOfClick"`
	URL            string         `json:"URL"`
	UserIP         string         `json:"UserIp"`
}

// URLClickAction .
type URLClickAction int

// MarshalJSON marshals into a string.
func (t URLClickAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// URLClickAction enum.
const (
	Blockpage URLClickAction = iota + 2
	PendingDetonationPage
	BlockPageOverride
	PendingDetonationPageOverride
)

// URLClickActionLiterals .
var URLClickActionLiterals = []string{
	"Blockpage",
	"PendingDetonationPage",
	"BlockPageOverride",
	"PendingDetonationPageOverride",
}

func (t URLClickAction) String() string {
	return URLClickActionLiterals[t]
}

// FileEvents .
type FileEvents struct {
	FileData         FileData       `json:"FileData"`
	SourceWorkload   SourceWorkload `json:"SourceWorkload"`
	DetectionMethod  string         `json:"DetectionMethod"`
	LastModifiedDate string         `json:"LastModifiedDate"`
	LastModifiedBy   string         `json:"LastModifiedBy"`
	EventDeepLink    string         `json:"EventDeepLink"`
}

// FileData .
type FileData struct {
	DocumentID    string      `json:"DocumentId"`
	FileName      string      `json:"FileName"`
	FilePath      string      `json:"FilePath"`
	FileVerdict   FileVerdict `json:"FileVerdict"`
	MalwareFamily string      `json:"MalwareFamily"`
	SHA256        string      `json:"SHA256"`
	FileSize      string      `json:"FileSize"`
}

// SourceWorkload .
type SourceWorkload int

// MarshalJSON marshals into a string.
func (t SourceWorkload) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// SourceWorkload enum.
const (
	SharePointOnlineWL SourceWorkload = iota
	OneDriveforBusinessWL
	MicrosoftTeamsWL
)

// SourceWorkloadLiterals .
var SourceWorkloadLiterals = []string{
	"SharePoint Online",
	"OneDrive for Business",
	"Microsoft Teams",
}

func (t SourceWorkload) String() string {
	return SourceWorkloadLiterals[t]
}
