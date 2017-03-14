package secrets

import "github.com/untoldwind/trustless/api"

type SecretBlock struct {
	ID      string             `json:"id"`
	Version *api.SecretVersion `json:"version"`
}
