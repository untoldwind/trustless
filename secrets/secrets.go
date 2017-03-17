package secrets

import "github.com/untoldwind/trustless/api"

type Secrets interface {
	IsInitialized() bool

	IsLocked() bool
	Lock()
	Unlock(name, email, passphrase string) error

	List() (*api.SecretList, error)
	Add(id string, secretType api.SecretType, version api.SecretVersion) error
	Get(secretID string) (*api.Secret, error)
}
