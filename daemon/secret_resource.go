package daemon

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/secrets"
)

type SecretResource struct {
	rest.ResourceBase
	logger   logging.Logger
	secrets  secrets.Secrets
	secretID string
}

func NewSecretResource(secrets secrets.Secrets, secretID string, logger logging.Logger) *SecretResource {
	return &SecretResource{
		logger:   logger.WithField("resource", "secret").WithField("secretID", secretID),
		secrets:  secrets,
		secretID: secretID,
	}
}

func (r *SecretResource) Self() rest.Link {
	return rest.SimpleLink(fmt.Sprintf("/v1/secrets/%s", url.QueryEscape(r.secretID)))
}

func (r *SecretResource) Get(request *http.Request) (interface{}, error) {
	secret, err := r.secrets.Get(r.secretID)

	if err == secrets.SecretNotFound {
		return nil, rest.NotFound
	} else if err != nil {
		return nil, err
	}
	return secret, nil
}
