package schema

import "encoding/json"

// Sway .
type Sway struct {
	ObjectType      ObjectType      `json:"ObjectType,omitempty"`
	Endpoint        Endpoint        `json:"Endpoint,omitempty"`
	BrowserName     string          `json:"BrowserName,omitempty"`
	DeviceType      DeviceType      `json:"DeviceType,omitempty"`
	SwayLookupID    string          `json:"SwayLookupId,omitempty"`
	SiteURL         string          `json:"SiteUrl,omitempty"`
	OperationResult OperationResult `json:"OperationResult,omitempty"`
}

// ObjectType  .
type ObjectType int

// MarshalJSON marshals into a string.
func (t ObjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// ObjectType  enum.
const (
	SwayType ObjectType = iota
	SwayEmbedded
	SwayAdminPortal
)

// ObjectTypeLiterals .
var ObjectTypeLiterals = []string{
	"Sway",
	"SwayEmbedded",
	"SwayAdminPortal",
}

func (t ObjectType) String() string {
	return ObjectTypeLiterals[t]
}

// OperationResult  .
type OperationResult int

// MarshalJSON marshals into a string.
func (t OperationResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// OperationResult  enum.
const (
	Succeeded OperationResult = iota
	Failed
)

// OperationResultLiterals .
var OperationResultLiterals = []string{
	"Succeeded",
	"Failed",
}

func (t OperationResult) String() string {
	return OperationResultLiterals[t]
}

// Endpoint  .
type Endpoint int

// MarshalJSON marshals into a string.
func (t Endpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// Endpoint  enum.
const (
	SwayWeb Endpoint = iota
	SwayIOS
	SwayWindows
	SwayAndroid
)

// EndpointLiterals .
var EndpointLiterals = []string{
	"SwayWeb",
	"SwayIOS",
	"SwayWindows",
	"SwayAndroid",
}

func (t Endpoint) String() string {
	return EndpointLiterals[t]
}

// DeviceType  .
type DeviceType int

// MarshalJSON marshals into a string.
func (t DeviceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// DeviceType  enum.
const (
	Desktop DeviceType = iota
	Mobile
	Tablet
)

// DeviceTypeLiterals .
var DeviceTypeLiterals = []string{
	"Desktop",
	"Mobile",
	"Tablet",
}

func (t DeviceType) String() string {
	return DeviceTypeLiterals[t]
}
