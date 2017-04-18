package pgp

import (
	"context"

	zxcvbn "github.com/nbutton23/zxcvbn-go"
	"github.com/untoldwind/trustless/api"
)

func (c *pgpSecrets) EstimateStrength(ctx context.Context, estimate api.PasswordEstimate) (*api.PasswordStrength, error) {
	result := zxcvbn.PasswordStrength(estimate.Password, estimate.Inputs)

	return &api.PasswordStrength{
		Entropy:          result.Entropy,
		CrackTime:        result.CrackTime,
		CrackTimeDisplay: result.CrackTimeDisplay,
		Score:            result.Score,
	}, nil
}
