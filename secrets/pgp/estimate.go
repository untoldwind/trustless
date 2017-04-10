package pgp

import (
	"context"

	zxcvbn "github.com/nbutton23/zxcvbn-go"
	"github.com/untoldwind/trustless/api"
)

func (c *pgpSecrets) EstimateStrength(ctx context.Context, password string, inputs []string) (*api.PasswordStrength, error) {
	result := zxcvbn.PasswordStrength(password, inputs)

	return &api.PasswordStrength{
		Entropy:          result.Entropy,
		CrackTime:        result.CrackTime,
		CrackTimeDisplay: result.CrackTimeDisplay,
		Score:            result.Score,
	}, nil
}
