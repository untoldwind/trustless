package pgp

import (
	"context"
	"errors"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

var estimateFields = []string{
	"password",
}

func (s *pgpSecrets) Get(ctx context.Context, secretID string) (*api.Secret, error) {
	if s.isLocked() {
		return nil, secrets.ErrSecretsLocked
	}
	s.logger.Infof("Get secret %s", secretID)

	if s.index == nil {
		if err := s.buildIndex(); err != nil {
			return nil, err
		}
	}
	entry, ok := s.index.Entries[secretID]
	if !ok {
		return nil, secrets.ErrSecretNotFound
	}
	if entry.ID != secretID {
		return nil, errors.New("Index integrety failure")
	}

	result := &api.Secret{
		SecretCurrent: api.SecretCurrent{
			ID:   entry.ID,
			Type: entry.Type,
		},
		PasswordStrengths: map[string]*api.PasswordStrength{},
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
		for _, estimateField := range estimateFields {
			if value, ok := result.Current.Properties[estimateField]; ok {
				inputs := result.Current.URLs
				inputs = append(inputs, result.Current.Name)
				passwordStrength, err := s.EstimateStrength(ctx, api.PasswordEstimate{
					Password: value,
					Inputs:   inputs,
				})
				if err != nil {
					s.logger.ErrorErr(err)
				}
				result.PasswordStrengths[estimateField] = passwordStrength
			}
		}
	}
	s.autolocker.Reset()
	return result, nil
}
