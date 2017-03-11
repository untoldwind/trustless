package api

// SecretType is the type of the secret to store
type SecretType string

const (
	// SecretTypeLogin is a generic login (of a webpage)
	SecretTypeLogin SecretType = "login"
	// SecretTypeNote is a general secret note
	SecretTypeNote SecretType = "note"
	// SecretTypeLicence is some kind of software licence (key)
	SecretTypeLicence SecretType = "licence"
)

// Secret holds all information of a secret
type Secret struct {
	Type       SecretType        `json:"type"`
	Properties map[string]string `json:"properties"`
}
