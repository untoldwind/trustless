package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cli "gopkg.in/urfave/cli.v2"

	"github.com/fatih/color"
	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/config"
	"github.com/untoldwind/trustless/secrets"
	"github.com/untoldwind/trustless/secrets/remote"
)

var boldRed = color.New(color.FgRed, color.Bold).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func withDetailedErrors(action cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		err := action(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", boldRed("ERROR"), err.Error())

			if GlobalFlags.Debug {
				fmt.Fprintln(os.Stderr)
				fmt.Fprintf(os.Stderr, "%+v\n", err)
			}

			return errors.New("")
		}
		return nil
	}
}

func createLogger() logging.Logger {
	loggingOptions := logging.Options{
		Backend:   "logrus",
		LogFile:   GlobalFlags.LogFile,
		LogFormat: GlobalFlags.LogFormat,
		Level:     logging.Info,
	}
	if GlobalFlags.Debug {
		loggingOptions.Level = logging.Debug
	}
	return logging.NewLogrusLogger(loggingOptions).
		WithContext(map[string]interface{}{"process": "trustless", "version": config.Version()})
}

func createRemote(logger logging.Logger) secrets.Secrets {
	return remote.NewRemoteSecrets(logger)
}

func createClientContext() context.Context {
	return context.Background()
}

func handleSignals(logger logging.Logger) error {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	for {
		sig := <-signalCh

		logger.Infof("Caught signal: %v", sig)

		shutdown := false
		if sig == os.Interrupt || sig == syscall.SIGTERM {
			shutdown = true
		}

		if shutdown {
			return nil
		}
	}
}
