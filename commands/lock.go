package commands

import cli "gopkg.in/urfave/cli.v2"

// LockCommand is commandline command to lock the store
var LockCommand = &cli.Command{
	Name:   "lock",
	Usage:  "Unlock the store",
	Action: withDetailedErrors(lockStore),
}

func lockStore(ctx *cli.Context) error {
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
