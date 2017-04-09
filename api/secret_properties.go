package api

import "sort"

type SecretProperty struct {
	Name      string
	Display   string
	MultiLine bool
	Blurred   bool
}

type SecretPropertyList []SecretProperty

func (p SecretPropertyList) Len() int           { return len(p) }
func (p SecretPropertyList) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p SecretPropertyList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SecretPropertyList) Sort()              { sort.Sort(p) }

var username = SecretProperty{Name: "username", Display: "Username"}
var password = SecretProperty{Name: "password", Display: "Password", Blurred: true}
var sid = SecretProperty{Name: "sid", Display: "Sid", Blurred: true}
var notes = SecretProperty{Name: "notes", Display: "Notes", MultiLine: true}
var regCode = SecretProperty{Name: "regCode", Display: "Licence code", MultiLine: true}
var regName = SecretProperty{Name: "regName", Display: "Licenced to"}
var productVersion = SecretProperty{Name: "productVersion", Display: "Version"}

var SecretProperties = SecretPropertyList{
	username,
	password,
	sid,
	notes,
	regName,
	regCode,
	productVersion,
}

type SecretTypeDefinition struct {
	Type       SecretType
	Display    string
	Properties []SecretProperty
}

var SecretTypes = []SecretTypeDefinition{
	{Type: SecretTypeLogin, Display: "Login", Properties: []SecretProperty{username, password, notes}},
	{Type: SecretTypeNote, Display: "Note", Properties: []SecretProperty{notes}},
	{Type: SecretTypeWLAN, Display: "WLAN", Properties: []SecretProperty{sid, password, notes}},
	{Type: SecretTypeLicence, Display: "Licence", Properties: []SecretProperty{regName, regCode, productVersion, notes}},
}
