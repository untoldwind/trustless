package remote

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

// Identities retrieves all known identities that can access the secret store
func (c *remoteSecrets) Identities(ctx context.Context) ([]api.Identity, error) {
	result, err := c.get(ctx, "/v1/identities")
	if err != nil {
		return nil, err
	}
	var identities []api.Identity
	if err := json.Unmarshal(result, &identities); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal identities")
	}
	return identities, nil
}
