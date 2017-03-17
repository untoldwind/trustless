package pgp

import (
	"errors"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

func (s *pgpSecrets) Get(secretID string) (*api.Secret, error) {
	if s.index == nil {
		if err := s.buildIndex(); err != nil {
			return nil, err
		}
	}
	entry, ok := s.index.Entries[secretID]
	if !ok {
		return nil, secrets.SecretNotFound
	}
	if entry.ID != secretID {
		return nil, errors.New("Index integrety failure")
	}

	result := &api.Secret{
		SecretCurrent: api.SecretCurrent{
			ID:   entry.ID,
			Type: entry.Type,
		},
	}

	for blockID := range entry.Blocks {
		block, err := s.store.GetBlock(blockID)
		if err != nil {
			return nil, err
		}
		if block == nil {
			continue
		}
		secretBlock, err := s.decryptSecret(block)
		if err != nil {
			return nil, err
		}
		if secretBlock.ID != secretID {
			return nil, errors.New("Index integrety failure")
		}
		result.Versions = append(result.Versions, &secretBlock.Version)
	}
	result.Versions.Sort()
	if len(result.Versions) > 0 {
		result.Current = result.Versions[0]
	}
	return result, nil
}
