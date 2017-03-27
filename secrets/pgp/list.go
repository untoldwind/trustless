package pgp

import (
	"context"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

func (s *pgpSecrets) List(ctx context.Context) (*api.SecretList, error) {
	if s.isLocked() {
		return nil, secrets.ErrSecretsLocked
	}
	s.logger.Info("List secrets")

	if err := s.buildIndex(); err != nil {
		return nil, err
	}

	s.autolocker.Reset()

	return s.index.list(), nil
}
