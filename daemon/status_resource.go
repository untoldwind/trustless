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
	locked, autolockAt := r.secrets.IsLocked()
	return &api.Status{
		Initialized: r.secrets.IsInitialized(),
		Locked:      locked,
		AutolockAt:  autolockAt,
		Version:     config.Version(),
	}, nil
}
