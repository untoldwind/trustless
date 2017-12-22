package cmds

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/untoldwind/trustless/daemon"
	"github.com/untoldwind/trustless/secrets/pgp"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start as daemon",
	Run:   withDetailedErrors(startDaemon),
}

func startDaemon(cmd *cobra.Command, args []string) error {
	logger := createLogger()

	secrets, err := pgp.NewPGPSecrets(cmdSettings.StoreURL, cmdSettings.NodeID, 4096, cmdSettings.UnlockTimeout, cmdSettings.UnlockTimeoutHard, logger)
	if err != nil {
		return err
	}
	defer secrets.Lock(context.Background())

	daemon := daemon.NewDaemon(secrets, logger)

	if err := daemon.Start(); err != nil {
		return err
	}
	defer daemon.Stop()

	return handleSignals(logger)
}
