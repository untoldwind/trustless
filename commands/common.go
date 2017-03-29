package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	cli "gopkg.in/urfave/cli.v2"
	yaml "gopkg.in/yaml.v2"

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

func readConfig(logger logging.Logger) (*config.CommonConfig, error) {
	raw, err := ioutil.ReadFile(GlobalFlags.ConfigFile)

	if os.IsNotExist(err) {
		logger.Warnf("Configuration %s does not exists, create default")
		defaultConfig, err := config.DefaultCommonConfig()
		if err != nil {
			return nil, err
		}
		if err := writeClientConfig(defaultConfig); err != nil {
			return nil, err
		}
		return defaultConfig, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "read config file failed")
	}

	var config config.CommonConfig
	if strings.HasSuffix(GlobalFlags.ConfigFile, ".json") {
		if err := json.Unmarshal(raw, &config); err != nil {
			return nil, errors.Wrap(err, "config.ClientConfig json unmarshal failed")
		}
	} else {
		if err := yaml.Unmarshal(raw, &config); err != nil {
			return nil, errors.Wrap(err, "config.ClientConfig yaml unmarshal failed")
		}
	}
	return &config, nil
}

func writeClientConfig(config *config.CommonConfig) error {
	var raw []byte
	var err error
	if strings.HasSuffix(GlobalFlags.ConfigFile, ".json") {
		raw, err = json.Marshal(config)
		if err != nil {
			return errors.Wrap(err, "config.CommonConfig json marshal failed")
		}
	} else {
		raw, err = yaml.Marshal(config)
		if err != nil {
			return errors.Wrap(err, "config.CommonConfig yaml marshal failed")
		}
	}
	if err := os.MkdirAll(filepath.Dir(GlobalFlags.ConfigFile), 0700); err != nil {
		return errors.Wrap(err, "creating config file directory failed")
	}
	if err := ioutil.WriteFile(GlobalFlags.ConfigFile, raw, 0600); err != nil {
		return errors.Wrap(err, "write config file failed")
	}
	return nil
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
