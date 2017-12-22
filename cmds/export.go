package cmds

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

var exportFlags = struct {
	format string
}{}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export store",
	Args:  cobra.MaximumNArgs(1),
	Run:   withDetailedErrors(exportFile),
}

func exportFile(cmd *cobra.Command, args []string) error {
	logger := createLogger()
	client := createRemote(logger)

	if _, err := unlockStore(client); err != nil {
		return err
	}

	fileName := ""
	if len(args) > 0 {
		fileName = args[0]
	}

	switch exportFlags.format {
	case "json":
		return exportJson(client, fileName)
	case "keepass":
		return exportKeepass(client, fileName)
	default:
		return errors.Errorf("Invalid export format: %s", exportFlags.format)
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
