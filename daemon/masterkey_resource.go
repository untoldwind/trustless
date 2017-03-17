package daemon

import (
	"encoding/json"
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

type MasterKeyResource struct {
	rest.ResourceBase
	logger  logging.Logger
	secrets secrets.Secrets
}

func NewMasterKeyResource(secrets secrets.Secrets, logger logging.Logger) *MasterKeyResource {
	return &MasterKeyResource{
		logger:  logger.WithField("resource", "masterkey"),
		secrets: secrets,
	}
}

func (MasterKeyResource) Self() rest.Link {
	return rest.SimpleLink("/v1/masterkey")
}

func (r *MasterKeyResource) Get(request *http.Request) (interface{}, error) {
	return &api.MasterKey{
		Locked: r.secrets.IsLocked(),
	}, nil
}

func (r *MasterKeyResource) Update(request *http.Request) (interface{}, error) {
	var unlock api.MasterKeyUnlock

	defer request.Body.Close()
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&unlock); err != nil {
		return nil, rest.BadRequest.WithDetails(err.Error())
	}
	if err := r.secrets.Unlock(unlock.Name, unlock.Email, unlock.Passphrase); err != nil {
		return nil, rest.InternalServerError(err)
	}
	return nil, nil
}

func (r *MasterKeyResource) Delete(request *http.Request) (interface{}, error) {
	r.secrets.Lock()
	return nil, nil
}
