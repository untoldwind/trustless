package commands

import (
	"github.com/untoldwind/trustless/daemon"
	"github.com/untoldwind/trustless/secrets"
	cli "gopkg.in/urfave/cli.v2"
)

var DaemonCommang = &cli.Command{
	Name:   "daemon",
	Usage:  "Start daemon",
	Action: withDetailedErrors(startDaemon),
}

func startDaemon(ctx *cli.Context) error {
	logger := createLogger()

	config, err := readConfig(logger)
	if err != nil {
		return err
	}

	secrets, err := secrets.NewSecrets(config.StoreURL, config.NodeID, logger)
	if err != nil {
		return err
	}

	daemon := daemon.NewDaemon(secrets, logger)

	if err := daemon.Start(); err != nil {
		return err
	}
	defer daemon.Stop()

	return handleSignals(logger)
}
