package remote

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

func (c *remoteSecrets) EstimateStrength(ctx context.Context, password string, inputs []string) (*api.PasswordStrength, error) {
	estimate, err := json.Marshal(&api.PasswordEstimate{
		Password: password,
		Inputs:   inputs,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialized password estimate")
	}
	result, err := c.post(ctx, "/v1/estimate", estimate)
	if err != nil {
		return nil, err
	}
	var passwordStrength api.PasswordStrength
	if err := json.Unmarshal(result, &passwordStrength); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize status")
	}
	return &passwordStrength, nil
}
