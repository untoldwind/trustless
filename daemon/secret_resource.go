package daemon

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/secrets"
)

// SecretResource is a REST resource representing a specific secret
type SecretResource struct {
	rest.ResourceBase
	logger   logging.Logger
	secrets  secrets.Secrets
	secretID string
}

// NewSecretResource create a SecretResource
func NewSecretResource(secrets secrets.Secrets, secretID string, logger logging.Logger) *SecretResource {
	return &SecretResource{
		logger:   logger.WithField("resource", "secret").WithField("secretID", secretID),
		secrets:  secrets,
		secretID: secretID,
	}
}

// Self link of the resource
func (r *SecretResource) Self() rest.Link {
	return rest.SimpleLink(fmt.Sprintf("/v1/secrets/%s", url.QueryEscape(r.secretID)))
}

// Get the secret the resource represents
func (r *SecretResource) Get(request *http.Request) (interface{}, error) {
	secret, err := r.secrets.Get(request.Context(), r.secretID)

	if err == secrets.ErrSecretNotFound {
		return nil, rest.HTTPNotFound
	} else if err != nil {
		return nil, err
	}
	return secret, nil
}
