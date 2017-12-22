package cmds

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/untoldwind/trustless/api"
)

var estimateCmd = &cobra.Command{
	Use:   "estimate",
	Short: "Estimate password strength",
	Run:   withDetailedErrors(estimatePassword),
}

func estimatePassword(cmd *cobra.Command, args []string) error {
	logger := createLogger()
	client := createRemote(logger)

	password, err := readPassphrase("Estimate password: ")
	if err != nil {
		return err
	}

	result, err := client.EstimateStrength(context.Background(), api.PasswordEstimate{
		Password: password,
	})

	fmt.Printf("Entropy  : %f\n", result.Entropy)
	fmt.Printf("Cracktime: %f (%s)\n", result.CrackTime, result.CrackTimeDisplay)
	fmt.Printf("Scope    : %d\n", result.Score)

	return nil
}
