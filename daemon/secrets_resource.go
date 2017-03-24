package daemon

import (
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/secrets"
)

// SecretsResource is a REST resource representing all secrets in the store
type SecretsResource struct {
	rest.ResourcesBase
	logger  logging.Logger
	secrets secrets.Secrets
}

// NewSecretsResource create a new SecretsResource
func NewSecretsResource(secrets secrets.Secrets, logger logging.Logger) *SecretsResource {
	return &SecretsResource{
		logger:  logger.WithField("resource", "secrets"),
		secrets: secrets,
	}
}

// Self link to the resource
func (SecretsResource) Self() rest.Link {
	return rest.SimpleLink("/v1/secrets")
}

// List all secrets in the store
func (r *SecretsResource) List(*http.Request) (interface{}, error) {
	return r.secrets.List()
}

// FindById looks up a secret by its id
func (r *SecretsResource) FindById(id string) (interface{}, error) {
	return NewSecretResource(r.secrets, id, r.logger), nil
}
