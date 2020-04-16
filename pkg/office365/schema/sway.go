package schema

import "encoding/json"

// Sway .
type Sway struct {
	AuditRecord
	ObjectType      *ObjectType      `json:"ObjectType,omitempty"`
	Endpoint        *Endpoint        `json:"Endpoint,omitempty"`
	BrowserName     *string          `json:"BrowserName,omitempty"`
	DeviceType      *DeviceType      `json:"DeviceType,omitempty"`
	SwayLookupID    *string          `json:"SwayLookupId,omitempty"`
	SiteURL         *string          `json:"SiteUrl,omitempty"`
	OperationResult *OperationResult `json:"OperationResult,omitempty"`
}

// ObjectType  .
type ObjectType int

// MarshalJSON marshals into a string.
func (t ObjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// ObjectType  enum.
const (
	SwayOT ObjectType = iota
	SwayEmbeddedOT
	SwayAdminPortalOT
)

func (t ObjectType) String() string {
	literals := map[ObjectType]string{
		SwayOT:            "Sway",
		SwayEmbeddedOT:    "SwayEmbedded",
		SwayAdminPortalOT: "SwayAdminPortal",
	}
	return literals[t]
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

func (t OperationResult) String() string {
	literals := map[OperationResult]string{
		Succeeded: "Succeeded",
		Failed:    "Failed",
	}
	return literals[t]
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

func (t Endpoint) String() string {
	literals := map[Endpoint]string{
		SwayWeb:     "SwayWeb",
		SwayIOS:     "SwayIOS",
		SwayWindows: "SwayWindows",
		SwayAndroid: "SwayAndroid",
	}
	return literals[t]
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

func (t DeviceType) String() string {
	literals := map[DeviceType]string{
		Desktop: "Desktop",
		Mobile:  "Mobile",
		Tablet:  "Tablet",
	}
	return literals[t]
}
