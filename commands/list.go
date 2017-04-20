package commands

import (
	"fmt"

	"github.com/untoldwind/trustless/api"

	cli "gopkg.in/urfave/cli.v2"
)

var ListCommand = &cli.Command{
	Name:   "list",
	Usage:  "List all secrets in the store",
	Action: withDetailedErrors(listSecrets),
}

func listSecrets(ctx *cli.Context) error {
	logger := createLogger()
	client := createRemote(logger)

	if _, err := unlockStore(client); err != nil {
		return err
	}

	secrets, err := client.List(createClientContext(), api.SecretListFilter{})
	if err != nil {
		return err
	}
	for _, entry := range secrets.Entries {
		fmt.Printf("%v\n", entry)
	}

	return nil
}
