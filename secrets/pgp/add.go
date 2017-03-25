package pgp

import (
	"context"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
	"github.com/untoldwind/trustless/store/model"
)

func (s *pgpSecrets) Add(ctx context.Context, id string, secretType api.SecretType, version api.SecretVersion) error {
	s.logger.Info("Add secret %s", id)

	secretBlock := &secrets.SecretBlock{
		ID:      id,
		Type:    secretType,
		Version: version,
	}
	encrypted, err := s.encryptSecret(secretBlock)
	if err != nil {
		return err
	}
	blockID, err := s.store.AddBlock(encrypted)
	if err != nil {
		return err
	}
	commitID, err := s.store.Commit(s.nodeID, []model.Change{
		{Operation: model.ChangeOpAdd, BlockID: blockID},
	})
	if s.index != nil {
		s.index.registerCommit(commitID, map[string]*secrets.SecretBlock{
			blockID: secretBlock,
		})
	}
	s.autolocker.Reset()
	return err
}
