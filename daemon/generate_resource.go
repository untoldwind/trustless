package daemon

import (
	"encoding/json"
	"net/http"

	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

type GenerateResource struct {
	rest.ResourceBase
	secrets secrets.Secrets
}

func NewGenerateResource(secrets secrets.Secrets) *GenerateResource {
	return &GenerateResource{
		secrets: secrets,
	}
}

func (GenerateResource) Self() rest.Link {
	return rest.SimpleLink("/v1/generate")
}

func (r *GenerateResource) Post(request *http.Request) (interface{}, error) {
	var parameter api.GenerateParameter

	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&parameter); err != nil {
		return nil, rest.HTTPBadRequest.WithDetails(err.Error())
	}

	return r.secrets.GeneratePassword(request.Context(), parameter)
}
