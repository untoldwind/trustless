package main

import (
	"fmt"
	"os"

	"github.com/untoldwind/trustless/commands"
	"github.com/untoldwind/trustless/config"

	cli "gopkg.in/urfave/cli.v2"
)

func showError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "%+v\n", err)
}

func main() {
	app := &cli.App{
		Name:    "trustless",
		Usage:   "Password storage (if you have less trust in all the others)",
		Version: config.Version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-file",
				Value:       "",
				Usage:       "Log to file instead stdout",
				EnvVars:     []string{"DEPLOY_CONTROL_LOG_FILE"},
				Destination: &commands.GlobalFlags.LogFile,
			},
			&cli.StringFlag{
				Name:        "log-format",
				Value:       "text",
				Usage:       "Log format to use (test, json, logstash)",
				EnvVars:     []string{"DEPLOY_CONTROL_LOG_FORMAT"},
				Destination: &commands.GlobalFlags.LogFormat,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Enable debug logging",
				EnvVars:     []string{"DEPLOY_CONTROL_DEBUG"},
				Destination: &commands.GlobalFlags.Debug,
			},
			&cli.StringFlag{
				Name:        "config-file",
				Usage:       "Client configuration file",
				Value:       config.DefaultConfigFile(),
				Destination: &commands.GlobalFlags.ConfigFile,
			},
		},
		Commands: []*cli.Command{
			commands.InfoCommand,
			commands.LockCommand,
			commands.UnlockCommand,
			commands.ListCommand,
			commands.ImportCommand,
			commands.DaemonCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		showError(err)
	}
}
