package api

type SecretProperty struct {
	Name      string
	Display   string
	MultiLine bool
	Blurred   bool
}

var username = SecretProperty{Name: "username", Display: "Username"}
var password = SecretProperty{Name: "password", Display: "Password", Blurred: true}
var sid = SecretProperty{Name: "sid", Display: "Sid", Blurred: true}
var notes = SecretProperty{Name: "notes", Display: "Notes", MultiLine: true}
var regCode = SecretProperty{Name: "regCode", Display: "Licence code", MultiLine: true}
var regName = SecretProperty{Name: "regName", Display: "Licenced to"}
var productVersion = SecretProperty{Name: "productVersion", Display: "Version"}

var SecretProperties = []SecretProperty{
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
