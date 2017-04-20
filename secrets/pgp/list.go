package pgp

import (
	"context"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

func (s *pgpSecrets) List(ctx context.Context, filter api.SecretListFilter) (*api.SecretList, error) {
	if s.isLocked() {
		return nil, secrets.ErrSecretsLocked
	}
	s.logger.Info("List secrets")

	if err := s.buildIndex(); err != nil {
		return nil, err
	}

	s.autolocker.Reset()

	list := s.index.list()

	filtered := make([]*api.SecretEntry, 0, len(list.Entries))
	for _, entry := range list.Entries {
		if entry.Matches(filter) {
			filtered = append(filtered, entry)
		}
	}
	return &api.SecretList{
		AllTags: list.AllTags,
		Entries: filtered,
	}, nil
}
