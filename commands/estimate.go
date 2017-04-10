package commands

import (
	"context"
	"fmt"

	cli "gopkg.in/urfave/cli.v2"
)

var EstimateCommand = &cli.Command{
	Name:   "estimate",
	Usage:  "Estimate strength of a password",
	Action: withDetailedErrors(estimatePassword),
}

func estimatePassword(ctx *cli.Context) error {
	logger := createLogger()
	client := createRemote(logger)

	password, err := readPassphrase("Estimate password: ")
	if err != nil {
		return err
	}

	result, err := client.EstimateStrength(context.Background(), password, nil)

	fmt.Printf("Entropy  : %f\n", result.Entropy)
	fmt.Printf("Cracktime: %f (%s)\n", result.CrackTime, result.CrackTimeDisplay)
	fmt.Printf("Scope    : %d\n", result.Score)

	return nil
}
