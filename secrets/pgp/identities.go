package pgp

import (
	"context"

	"github.com/untoldwind/trustless/api"
)

func (s *pgpSecrets) Identities(ctx context.Context) ([]api.Identity, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	result := make([]api.Identity, 0, len(s.entities))
	for _, entity := range s.entities {
		if entity.PrivateKey == nil {
			continue
		}
		for _, identity := range entity.Identities {
			result = append(result, api.Identity{
				Name:  identity.UserId.Name,
				Email: identity.UserId.Email,
			})
		}
	}
	return result, nil
}
