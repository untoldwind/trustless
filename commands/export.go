package commands

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/tobischo/gokeepasslib"

	"github.com/pkg/errors"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
	cli "gopkg.in/urfave/cli.v2"
)

var ExportFlags = struct {
	Format string
}{}

var ExportCommand = &cli.Command{
	Name:      "export",
	Usage:     "Export store",
	ArgsUsage: "[<file>]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "format",
			Usage:       "Export format (json, keepass)",
			Value:       "json",
			Destination: &ExportFlags.Format,
		},
	},
	Action: withDetailedErrors(exportFile),
}

func exportFile(ctx *cli.Context) error {
	logger := createLogger()
	client := createRemote(logger)

	if _, err := unlockStore(client); err != nil {
		return err
	}

	fileName := ""
	if ctx.Args().Len() > 0 {
		fileName = ctx.Args().Get(0)
	}

	switch ExportFlags.Format {
	case "json":
		return exportJson(client, fileName)
	case "keepass":
		return exportKeepass(client, fileName)
	default:
		return errors.Errorf("Invalid export format: %s", ExportFlags.Format)
	}
}

func exportJson(client secrets.Secrets, fileName string) error {
	secrets, err := client.List(createClientContext(), api.SecretListFilter{})
	if err != nil {
		return err
	}

	var out io.Writer = os.Stdout

	if fileName != "" {
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return errors.Wrapf(err, "Unable to open: %s", fileName)
		}
		defer file.Close()

		out = file
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

func exportKeepass(client secrets.Secrets, fileName string) error {
	database := gokeepasslib.NewDatabase()

	root := gokeepasslib.NewGroup()
	root.Name = "trustless"

	secrets, err := client.List(createClientContext(), api.SecretListFilter{})
	if err != nil {
		return err
	}
	for _, entry := range secrets.Entries {
		secret, err := client.Get(createClientContext(), entry.ID)
		if err != nil {
			return err
		}
		entry := gokeepasslib.NewEntry()
		entry.Values = []gokeepasslib.ValueData{
			{Key: "Title", Value: gokeepasslib.V{Content: secret.Current.Name}},
		}
		entry.Tags = strings.Join(secret.Current.Tags, ",")
		root.Entries = append(root.Entries, entry)
	}

	database.Content = &gokeepasslib.DBContent{
		Meta: gokeepasslib.NewMetaData(),
		Root: &gokeepasslib.RootData{
			Groups: []gokeepasslib.Group{root},
		},
	}
	database.LockProtectedEntries()

	var out io.Writer = os.Stdout

	if fileName != "" {
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return errors.Wrapf(err, "Unable to open: %s", fileName)
		}
		defer file.Close()

		out = file
	}

	encoder := gokeepasslib.NewEncoder(out)

	return encoder.Encode(database)
}
