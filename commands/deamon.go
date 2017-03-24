package commands

import (
	"github.com/untoldwind/trustless/daemon"
	"github.com/untoldwind/trustless/secrets/pgp"
	cli "gopkg.in/urfave/cli.v2"
)

// DaemonCommand is the command line command to start the trustless daemon
var DaemonCommand = &cli.Command{
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

	secrets, err := pgp.NewPGPSecrets(config.StoreURL, config.NodeID, 4096, logger)
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
