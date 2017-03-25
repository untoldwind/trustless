package remote

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

func (c *remoteSecrets) List(ctx context.Context) (*api.SecretList, error) {
	result, err := c.get(ctx, "/v1/secrets")
	if err != nil {
		return nil, err
	}
	var list api.SecretList
	if err := json.Unmarshal(result, &list); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize status")
	}
	return &list, nil
}
