package main

import (
	"io"
	"log/syslog"
	"os"

	"github.com/leanovate/microtools/logging"
	"github.com/untoldwind/trustless/config"
)

func createLogger() logging.Logger {
	var output io.Writer
	output, err := syslog.New(syslog.LOG_LOCAL6, "trustless")
	if err != nil {
		output = os.Stdout
	}
	loggingOptions := logging.Options{
		Backend: "logrus",
		Output:  output,
		Level:   logging.Info,
	}
	return logging.NewLogrusLogger(loggingOptions).
		WithContext(map[string]interface{}{"process": "trustless-native", "version": config.Version()})
}
