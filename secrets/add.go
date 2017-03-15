package secrets

import (
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/store/model"
)

func (s *Secrets) Add(id string, secretType api.SecretType, version api.SecretVersion) error {
	secretBlock := &SecretBlock{
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
		s.index.registerCommit(commitID, map[string]*SecretBlock{
			blockID: secretBlock,
		})
	}
	return err
}
