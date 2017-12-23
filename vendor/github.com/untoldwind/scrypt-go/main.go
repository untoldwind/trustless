package main

import (
	"fmt"
	"os"

	"github.com/untoldwind/scrypt-go/commands"
	"github.com/untoldwind/scrypt-go/config"
	cli "gopkg.in/urfave/cli.v2"
)

func showError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "%+v\n", err)
}

func main() {
	app := &cli.App{
		Name:    "scrypt-go",
		Usage:   "Encrypt/decrypt files using scrypted PSK",
		Version: config.Version(),
		Commands: []*cli.Command{
			commands.EncryptCommand,
			commands.DecryptCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		showError(err)
	}
}
