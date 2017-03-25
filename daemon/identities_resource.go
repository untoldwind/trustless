package daemon

import (
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/secrets"
)

// IdentitiesResource is a REST resource representing all identities of the store
type IdentitiesResource struct {
	rest.ResourcesBase
	logger  logging.Logger
	secrets secrets.Secrets
}

// NewIdentitiesResource create a new IdentitiesResource
func NewIdentitiesResource(secrets secrets.Secrets, logger logging.Logger) *IdentitiesResource {
	return &IdentitiesResource{
		logger:  logger.WithField("resource", "secrets"),
		secrets: secrets,
	}
}

// Self link to the resource
func (IdentitiesResource) Self() rest.Link {
	return rest.SimpleLink("/v1/identities")
}

// List all identities in the store
func (r *IdentitiesResource) List(request *http.Request) (interface{}, error) {
	return r.secrets.Identities(request.Context())
}
