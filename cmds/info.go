package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/untoldwind/trustless/config"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get current status information",
	Run:   withDetailedErrors(getStatusInfo),
}

func getStatusInfo(cmd *cobra.Command, args []string) error {
	logger := createLogger()
	client := createRemote(logger)

	status, err := client.Status(createClientContext())
	if err != nil {
		return err
	}

	fmt.Printf("Client version is: %s\n", yellow(config.Version()))
	fmt.Printf("Daemon version is: %s\n", yellow(status.Version))
	reportStatus(status)

	return nil
}
