package commands

import (
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/untoldwind/trustless/api"
	cli "gopkg.in/urfave/cli.v2"
)

var ExportCommand = &cli.Command{
	Name:      "export",
	Usage:     "Export store",
	ArgsUsage: "[<file>]",
	Action:    withDetailedErrors(exportFile),
}

func exportFile(ctx *cli.Context) error {
	logger := createLogger()
	client := createRemote(logger)

	if _, err := unlockStore(client); err != nil {
		return err
	}

	var out io.Writer = os.Stdout

	if ctx.Args().Len() > 0 {
		file, err := os.OpenFile(ctx.Args().Get(0), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return errors.Wrapf(err, "Unable to open: %s", ctx.Args().Get(0))
		}
		defer file.Close()

		out = file
	}

	secrets, err := client.List(createClientContext(), api.SecretListFilter{})
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(out)
	for _, entry := range secrets.Entries {
		secret, err := client.Get(createClientContext(), entry.ID)
		if err != nil {
			return err
		}
		if err := encoder.Encode(secret); err != nil {
			return err
		}
	}

	return nil
}
