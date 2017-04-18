package daemon

import (
	"encoding/json"
	"net/http"

	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

type EstimateResource struct {
	rest.ResourceBase
	secrets secrets.Secrets
}

func NewEstimateResource(secrets secrets.Secrets) *EstimateResource {
	return &EstimateResource{
		secrets: secrets,
	}
}

func (EstimateResource) Self() rest.Link {
	return rest.SimpleLink("/v1/estimate")
}

func (r *EstimateResource) Post(request *http.Request) (interface{}, error) {
	var estimate api.PasswordEstimate

	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&estimate); err != nil {
		return nil, rest.HTTPBadRequest.WithDetails(err.Error())
	}

	return r.secrets.EstimateStrength(request.Context(), estimate)
}
