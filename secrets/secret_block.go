package secrets

import "github.com/untoldwind/trustless/api"

type SecretBlock struct {
	ID      string            `json:"id"`
	Type    api.SecretType    `json:"type"`
	Version api.SecretVersion `json:"version"`
}