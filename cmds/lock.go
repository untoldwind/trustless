package cmds

import (
	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the store",
	Run:   withDetailedErrors(lockStore),
}

func lockStore(cmd *cobra.Command, args []string) error {
	logger := createLogger()
	client := createRemote(logger)

	if err := client.Lock(createClientContext()); err != nil {
		return err
	}

	status, err := client.Status(createClientContext())
	if err != nil {
		return err
	}
	reportStatus(status)

	return nil
}
