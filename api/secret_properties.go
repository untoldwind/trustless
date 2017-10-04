package api

import "sort"

type SecretProperty struct {
	Name      string
	Display   string
	MultiLine bool
	Blurred   bool
	OTPParams bool
}

type SecretPropertyList []SecretProperty

func (p SecretPropertyList) Len() int           { return len(p) }
func (p SecretPropertyList) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p SecretPropertyList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SecretPropertyList) Sort()              { sort.Sort(p) }

var PropertyUsername = SecretProperty{Name: "username", Display: "Username"}
var PropertyPassword = SecretProperty{Name: "password", Display: "Password", Blurred: true}
var PropertySID = SecretProperty{Name: "sid", Display: "Sid", Blurred: true}
var PropertyNotes = SecretProperty{Name: "notes", Display: "Notes", MultiLine: true}
var PropertyRegCode = SecretProperty{Name: "regCode", Display: "Licence code", MultiLine: true}
var PropertyRegName = SecretProperty{Name: "regName", Display: "Licenced to"}
var PropertyProductVersion = SecretProperty{Name: "productVersion", Display: "Version"}
var PropertyTOTPUrl = SecretProperty{Name: "totpUrl", Display: "TOTP Url", OTPParams: true}

var SecretProperties = SecretPropertyList{
	PropertyUsername,
	PropertyPassword,
	PropertyTOTPUrl,
	PropertySID,
	PropertyNotes,
	PropertyRegName,
	PropertyRegCode,
	PropertyProductVersion,
}

type SecretTypeDefinition struct {
	Type       SecretType
	Display    string
	Properties []SecretProperty
}

var SecretTypes = []SecretTypeDefinition{
	{Type: SecretTypeLogin, Display: "Login", Properties: []SecretProperty{
		PropertyUsername,
		PropertyPassword,
		PropertyTOTPUrl,
		PropertyNotes,
	}},
	{Type: SecretTypeNote, Display: "Note", Properties: []SecretProperty{
		PropertyNotes,
	}},
	{Type: SecretTypeWLAN, Display: "WLAN", Properties: []SecretProperty{
		PropertySID,
		PropertyPassword,
		PropertyNotes,
	}},
	{Type: SecretTypeLicence, Display: "Licence", Properties: []SecretProperty{
		PropertyRegName,
		PropertyRegCode,
		PropertyProductVersion,
		PropertyNotes,
	}},
}
