package pgp

import (
	"context"

	"golang.org/x/crypto/openpgp"

	"github.com/untoldwind/trustless/api"
)

func (s *pgpSecrets) Identities(ctx context.Context) ([]api.Identity, error) {
	return s.identities, nil
}

func identitiesFromEntities(entities openpgp.EntityList) []api.Identity {
	result := make([]api.Identity, 0, len(entities))
	for _, entity := range entities {
		for _, identity := range entity.Identities {
			result = append(result, api.Identity{
				Name:  identity.UserId.Name,
				Email: identity.UserId.Email,
			})
		}
	}
	return result
}
