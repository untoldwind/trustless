package remote

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

func (c *remoteSecrets) Get(ctx context.Context, secretID string) (*api.Secret, error) {
	result, err := c.get(ctx, "/v1/secrets/"+url.QueryEscape(secretID))
	if err != nil {
		return nil, err
	}
	var secret api.Secret
	if err := json.Unmarshal(result, &secret); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize status")
	}
	return &secret, nil
}
