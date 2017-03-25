package commands

import (
	"fmt"

	"github.com/untoldwind/trustless/config"
	cli "gopkg.in/urfave/cli.v2"
)

// InfoCommand is the commandline command to retrieve all relevant status information
var InfoCommand = &cli.Command{
	Name:   "info",
	Usage:  "Get current status information",
	Action: withDetailedErrors(getStatusInfo),
}

func getStatusInfo(ctx *cli.Context) error {
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
