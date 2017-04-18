package pgp

import (
	"context"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets/generate"
)

func (c *pgpSecrets) GeneratePassword(ctx context.Context, parameter api.GenerateParameter) (string, error) {
	return generate.Password(parameter)
}
