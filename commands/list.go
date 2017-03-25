package commands

import (
	"fmt"

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

	secrets, err := client.List(createClientContext())
	if err != nil {
		return err
	}
	for _, entry := range secrets.Entries {
		fmt.Printf("%v\n", entry)
	}

	return nil
}
