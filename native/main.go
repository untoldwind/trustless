package main

import (
	"os"

	"github.com/leanovate/microtools/logging"
	"github.com/untoldwind/trustless/secrets"
	"github.com/untoldwind/trustless/secrets/remote"
)

func createRemote(logger logging.Logger) secrets.Secrets {
	return remote.NewRemoteSecrets(logger)
}

func main() {
	logger := createLogger()
	secrets := createRemote(logger)
	logger.Info("Started trustless-native")
	for {
		command, err := readCommand(os.Stdin)
		if err != nil {
			logger.ErrorErr(err)
			os.Exit(1)
		}
		if command == nil {
			os.Exit(0)
		}
		reply, err := process(command, secrets)
		if err != nil {
			logger.ErrorErr(err)
			os.Exit(1)
		}
		if err := writeReply(os.Stdout, command.Command, reply); err != nil {
			logger.ErrorErr(err)
			os.Exit(1)
		}
	}
}
