package daemon

import (
	"net/http"

	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/config"
	"github.com/untoldwind/trustless/secrets"
)

type StatusResource struct {
	rest.ResourceBase
	secrets secrets.Secrets
}

func NewStatusResource(secrets secrets.Secrets) *StatusResource {
	return &StatusResource{
		secrets: secrets,
	}
}

func (StatusResource) Self() rest.Link {
	return rest.SimpleLink("/status")
}

func (r *StatusResource) Get(request *http.Request) (interface{}, error) {
	return &api.Status{
		Initialized: r.secrets.IsInitialized(),
		Locked:      r.secrets.IsLocked(),
		Version:     config.Version(),
	}, nil
}
