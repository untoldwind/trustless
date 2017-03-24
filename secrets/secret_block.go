package secrets

import "github.com/untoldwind/trustless/api"

// SecretBlock represents a block containing the version of a secret.
// This is what actually has to be encrypted and stored to the underlying store
// implementation.
type SecretBlock struct {
	ID      string            `json:"id"`
	Type    api.SecretType    `json:"type"`
	Version api.SecretVersion `json:"version"`
}
