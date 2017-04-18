package remote

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

func (c *remoteSecrets) GeneratePassword(ctx context.Context, parameter api.GenerateParameter) (string, error) {
	param, err := json.Marshal(&parameter)
	if err != nil {
		return "", errors.Wrap(err, "Failed to serialized generate parameter")
	}
	result, err := c.post(ctx, "/v1/generate", param)
	if err != nil {
		return "", err
	}
	var password string
	if err := json.Unmarshal(result, &password); err != nil {
		return "", errors.Wrap(err, "Failed to deserialize generated password")
	}
	return password, nil
}
