package commands

import cli "gopkg.in/urfave/cli.v2"

var DaemonCommang = &cli.Command{
	Name:   "daemon",
	Usage:  "Start daemon",
	Action: withDetailedErrors(startDaemon),
}

func startDaemon(ctx *cli.Context) error {
	return nil
}
