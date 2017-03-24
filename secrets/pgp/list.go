package pgp

import (
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

func (s *pgpSecrets) List() (*api.SecretList, error) {
	if s.IsLocked() {
		return nil, secrets.ErrSecretsLocked
	}
	if err := s.buildIndex(); err != nil {
		return nil, err
	}

	return s.index.list(), nil
}
