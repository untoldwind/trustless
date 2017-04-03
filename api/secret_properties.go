package api

type SecretProperty struct {
	Name    string
	Display string
	Hidden  bool
	Blurred bool
}

var SecretProperties = []SecretProperty{
	{Name: "username", Display: "Username"},
	{Name: "password", Display: "Password", Blurred: true},
	{Name: "notes", Display: "Notes"},
	{Name: "sid", Display: "Sid"},
}
