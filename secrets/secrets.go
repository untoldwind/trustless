package secrets

import (
	"context"

	"github.com/untoldwind/trustless/api"
)

// Secrets is the interface to any secret store implementation
type Secrets interface {
	// Status gets the current status of the store.
	Status(ctx context.Context) (*api.Status, error)
	// Lock the store
	Lock(ctx context.Context) error
	// Unlock the store for a given identity
	Unlock(ctx context.Context, name, email, passphrase string) error

	// List all identities that have access to the store
	Identities(ctx context.Context) ([]api.Identity, error)

	// List all secrets of the store (only references)
	List(ctx context.Context) (*api.SecretList, error)
	// Add a secret to the store
	Add(ctx context.Context, id string, secretType api.SecretType, version api.SecretVersion) error
	// Get a secret from the store
	Get(ctx context.Context, secretID string) (*api.Secret, error)

	// EstimateStrength of a passwrd
	EstimateStrength(ctx context.Context, password string, inputs []string) (*api.PasswordStrength, error)
}
