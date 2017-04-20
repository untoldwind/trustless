package daemon

import (
	"encoding/json"
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/api"
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

// Create adds a secret to the store
func (r *SecretsResource) Create(request *http.Request) (rest.Resource, error) {
	var secretCurrent api.SecretCurrent

	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&secretCurrent); err != nil {
		return nil, rest.HTTPBadRequest.WithDetails(err.Error())
	}
	if secretCurrent.Current == nil {
		return nil, rest.HTTPBadRequest.WithDetails("No current secret")
	}

	if err := r.secrets.Add(request.Context(), secretCurrent.ID, secretCurrent.Type, *secretCurrent.Current); err != nil {
		return nil, err
	}
	return NewSecretResource(r.secrets, secretCurrent.ID, r.logger), nil
}

// List all secrets in the store
func (r *SecretsResource) List(request *http.Request) (interface{}, error) {
	return r.secrets.List(request.Context(), api.SecretListFilter{
		URL:  request.FormValue("url"),
		Tag:  request.FormValue("tag"),
		Type: api.SecretType(request.FormValue("type")),
		Name: request.FormValue("name"),
	})
}

// FindById looks up a secret by its id
func (r *SecretsResource) FindById(id string) (interface{}, error) {
	return NewSecretResource(r.secrets, id, r.logger), nil
}
