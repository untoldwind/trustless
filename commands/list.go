package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	fzf "github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/tui"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/config"
	"github.com/untoldwind/trustless/secrets"

	cli "gopkg.in/urfave/cli.v2"
)

var ListCommand = &cli.Command{
	Name:   "list",
	Usage:  "List all secrets in the store",
	Action: withDetailedErrors(listSecrets),
}

type secretsCommand struct {
	secrets  secrets.Secrets
	property string
}

func (s *secretsCommand) GetPreview(stripAnsi bool, delimiter fzf.Delimiter, query string, allItems []*fzf.Item) string {
	if len(allItems) == 0 {
		return ""
	}
	secretID := strings.Split(allItems[0].AsString(false), "\000")[0]
	secret, err := s.secrets.Get(context.Background(), secretID)
	if err != nil {
		return err.Error()
	}
	var lines []string
	lines = append(lines, fmt.Sprintf("Name        : %s", secret.Current.Name))
	lines = append(lines, fmt.Sprintf("Timestamp   : %s", secret.Current.Timestamp.String()))
	for _, url := range secret.Current.URLs {
		lines = append(lines, fmt.Sprintf("URL         : %s", url))
	}
	for name, value := range secret.Current.Properties {
		lines = append(lines, fmt.Sprintf("%-12s: %s", name, value))

	}
	return strings.Join(lines, "\n")
}

func (s *secretsCommand) HasPlusFlag() bool {
	return false
}

func (s *secretsCommand) Execute(withStdio bool, stripAnsi bool, delimiter fzf.Delimiter, forcePlus bool, query string, allItems []*fzf.Item) {
	if !withStdio {
		return
	}
	if len(allItems) == 0 {
		return
	}
	secretID := strings.Split(allItems[0].AsString(false), "\000")[0]
	secret, err := s.secrets.Get(context.Background(), secretID)
	if err != nil {

	}
	fmt.Fprintf(os.Stdout, secret.Current.Properties[s.property])
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

	opts := fzf.DefaultOptions()
	channel := make(chan []byte)
	opts.ReaderFactory = fzf.NewChannelReader(channel)
	delimter := "\000"
	opts.Delimiter.Str = &delimter
	opts.WithNth = []fzf.Range{{Begin: 2, End: 2}}
	opts.Preview.Command = &secretsCommand{secrets: client}
	opts.Printer = func(item string) {
		secretID := strings.Split(item, "\000")[0]
		secret, err := client.Get(context.Background(), secretID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, secret.Current.Properties["username"])
		fmt.Fprintln(os.Stdout, secret.Current.Properties["password"])
	}
	opts.Keymap[tui.CtrlSpace] = []fzf.Action{
		{Type: fzf.ActionTypeExecute, Command: &secretsCommand{secrets: client, property: "username"}},
		{Type: fzf.ActionTypeAbort},
	}
	opts.Keymap[tui.AltSpace] = []fzf.Action{
		{Type: fzf.ActionTypeExecute, Command: &secretsCommand{secrets: client, property: "password"}},
		{Type: fzf.ActionTypeAbort},
	}

	go func() {
		for _, entry := range secrets.Entries {
			channel <- []byte(entry.ID + "\000" + entry.Name)
		}
		close(channel)
	}()

	fzf.Run(opts, config.Version())

	return nil
}
