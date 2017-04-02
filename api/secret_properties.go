package api

type SecretProperty struct {
	Display string
	Hidden  bool
	Blurred bool
}

var SecretProperties = map[string]SecretProperty{
	"username": {Display: "Username"},
	"password": {Display: "Password", Blurred: true},
}
