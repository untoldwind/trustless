package pgp

import (
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

func (s *pgpSecrets) List() (*api.SecretList, error) {
	s.logger.Info("List secrets")

	if locked, _ := s.IsLocked(); locked {
		return nil, secrets.ErrSecretsLocked
	}
	if err := s.buildIndex(); err != nil {
		return nil, err
	}

	s.autolocker.Reset()

	return s.index.list(), nil
}
