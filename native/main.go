package main

import (
	"os"
	"runtime"

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

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 8192)
			count := runtime.Stack(trace, true)
			logger.Error(err)
			logger.Error(string(trace[0:count]))
			os.Exit(1)
		}
	}()

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
		reply, commandErr := process(command, secrets)
		if commandErr != nil {
			logger.ErrorErr(commandErr)
		}
		if err := writeReply(os.Stdout, command.Command, reply, commandErr); err != nil {
			logger.ErrorErr(err)
			os.Exit(1)
		}
	}
}
