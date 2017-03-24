package client

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

// Lock the secret store
func (c *Client) Lock(ctx context.Context) error {
	_, err := c.delete(ctx, "/v1/masterkey")
	return err
}

// Unlock the secret store
func (c *Client) Unlock(ctx context.Context, unlock api.MasterKeyUnlock) error {
	unlockRequest, err := json.Marshal(unlock)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal unlock request")
	}
	_, err = c.put(ctx, "/v1/masterkey", unlockRequest)
	return err
}
