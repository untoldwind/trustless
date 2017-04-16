package main

import "github.com/leanovate/microtools/logging"

func createLogger() logging.Logger {
	return logging.NewSimpleLoggerNull()
}
