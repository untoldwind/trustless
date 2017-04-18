package remote

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

func (c *remoteSecrets) EstimateStrength(ctx context.Context, estimate api.PasswordEstimate) (*api.PasswordStrength, error) {
	param, err := json.Marshal(&estimate)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialized password estimate")
	}
	result, err := c.post(ctx, "/v1/estimate", param)
	if err != nil {
		return nil, err
	}
	var passwordStrength api.PasswordStrength
	if err := json.Unmarshal(result, &passwordStrength); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize PasswordStrength")
	}
	return &passwordStrength, nil
}
