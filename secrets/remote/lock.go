package remote

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

// Status gets the current daemon status (and implicitly tests if the daemon is available)
func (c *remoteSecrets) Status(ctx context.Context) (*api.Status, error) {
	result, err := c.get(ctx, "/status")
	if err != nil {
		return nil, err
	}
	var status api.Status
	if err := json.Unmarshal(result, &status); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize status")
	}
	return &status, nil
}

// Lock the secret store
func (c *remoteSecrets) Lock(ctx context.Context) error {
	_, err := c.delete(ctx, "/v1/masterkey")
	return err
}

// Unlock the secret store
func (c *remoteSecrets) Unlock(ctx context.Context, name, email, passphrase string) error {
	unlockRequest, err := json.Marshal(api.MasterKeyUnlock{
		Identity: api.Identity{
			Name:  name,
			Email: email,
		},
		Passphrase: passphrase,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshal unlock request")
	}
	_, err = c.put(ctx, "/v1/masterkey", unlockRequest)
	return err
}
