package secrets

import "github.com/untoldwind/trustless/api"

// Secrets is the interface to any secret store implementation
type Secrets interface {
	// IsInitialized checks if the store has been initialized yet
	IsInitialized() bool

	// IsLocked checks if the the store is currently locked
	IsLocked() bool
	// Lock the store
	Lock()
	// Unlock the store for a given identity
	Unlock(name, email, passphrase string) error

	// List all identities that have access to the store
	Identities() ([]api.Identity, error)

	// List all secrets of the store (only references)
	List() (*api.SecretList, error)
	// Add a secret to the store
	Add(id string, secretType api.SecretType, version api.SecretVersion) error
	// Get a secret from the store
	Get(secretID string) (*api.Secret, error)
}
