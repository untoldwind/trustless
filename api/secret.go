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

type SecretCurrent struct {
	ID      string        `json:"id"`
	Type    SecretType    `json:"type"`
	Current SecretVersion `json:"current"`
}

// Secret holds all information of a secret
type Secret struct {
	SecretCurrent
	Version SecretVersions `json:"versions"`
}
