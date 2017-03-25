package remote

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

func (c *remoteSecrets) Add(ctx context.Context, id string, secretType api.SecretType, version api.SecretVersion) error {
	addRequest, err := json.Marshal(api.SecretCurrent{
		ID:      id,
		Type:    secretType,
		Current: &version,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal add secret request")
	}
	_, err = c.post(ctx, "/v1/secrets", addRequest)
	return err

}
