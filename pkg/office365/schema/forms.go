package schema

import "encoding/json"

// MicrosoftForms .
type MicrosoftForms struct {
	FormsUserTypes     []FormsUserTypes `json:"FormsUserTypes"`
	SourceApp          string           `json:"SourceApp"`
	FormName           string           `json:"FormName,omitempty"`
	FormID             string           `json:"FormId,omitempty"`
	FormTypes          []FormTypes      `json:"FormTypes,omitempty"`
	ActivityParameters string           `json:"ActivityParameters,omitempty"`
}

// FormsUserTypes .
type FormsUserTypes int

// MarshalJSON marshals into a string.
func (t FormsUserTypes) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// FormsUserTypes enum.
const (
	AdminUT FormsUserTypes = iota
	OwnerUT
	ResponderUT
	CoauthorUT
)

// FormsUserTypesLiterals .
var FormsUserTypesLiterals = []string{
	"Admin",
	"Owner",
	"Responder",
	"Coauthor",
}

func (t FormsUserTypes) String() string {
	return FormsUserTypesLiterals[t]
}

// FormTypes .
type FormTypes int

// MarshalJSON marshals into a string.
func (t FormTypes) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// FormTypes enum.
const (
	Form FormTypes = iota
	Quiz
	Survey
)

// FormTypesLiterals .
var FormTypesLiterals = []string{
	"Form",
	"Quiz",
	"Survey",
}

func (t FormTypes) String() string {
	return FormTypesLiterals[t]
}
