package daemon

import (
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/secrets"
)

type SecretsResource struct {
	rest.ResourcesBase
	logger  logging.Logger
	secrets secrets.Secrets
}

func NewSecretsResource(secrets secrets.Secrets, logger logging.Logger) *SecretsResource {
	return &SecretsResource{
		logger:  logger.WithField("resource", "secrets"),
		secrets: secrets,
	}
}

func (SecretsResource) Self() rest.Link {
	return rest.SimpleLink("/v1/secrets")
}

func (r *SecretsResource) List(*http.Request) (interface{}, error) {
	return r.secrets.List()
}

func (r *SecretsResource) FindById(id string) (interface{}, error) {
	return NewSecretResource(r.secrets, id, r.logger), nil
}
