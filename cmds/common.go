package cmds

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/awnumar/memguard"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/leanovate/microtools/logging"
	"github.com/spf13/cobra"
	"github.com/untoldwind/trustless/config"
	"github.com/untoldwind/trustless/secrets"
	"github.com/untoldwind/trustless/secrets/remote"
)

var boldRed = color.New(color.FgRed, color.Bold).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func showError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())

	if cmdSettings.Debug {
		fmt.Fprintln(os.Stderr)
		spew.Fdump(os.Stderr, err)
	}
	os.Exit(1)
}

func withDetailedErrors(action func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := action(cmd, args)
		if err != nil {
			showError(err)
			memguard.SafeExit(1)
		}
		memguard.SafeExit(0)
	}
}

func createLogger() logging.Logger {
	loggingOptions := logging.Options{
		Backend:   "simple",
		LogFile:   cmdSettings.LogFile,
		LogFormat: cmdSettings.LogFormat,
		Level:     logging.Info,
	}
	if cmdSettings.Debug {
		loggingOptions.Level = logging.Debug
	}
	return logging.NewLogger(loggingOptions).
		WithContext(map[string]interface{}{"process": "trustless", "version": config.Version()})
}

func createClientContext() context.Context {
	return context.Background()
}

func createRemote(logger logging.Logger) secrets.Secrets {
	return remote.NewRemoteSecrets(logger)
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
		fmt.Println("GC")
		runtime.GC()

		if shutdown {
			return nil
		}
	}
}
