package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"syscall"
)

// Logger is a generic logger interface to the different logger implementations
type Logger interface {
	// ErrorErr simply logs an error (inclduing stack if created vis pkg/errors)
	ErrorErr(error)
	// Errorf logs a formated error message
	Errorf(format string, args ...interface{})
	// Error logs an error message
	Error(args ...interface{})
	// Warnf logs a formatted warning message
	Warnf(format string, args ...interface{})
	// Warn logs a warning message
	Warn(args ...interface{})
	// Infof logs a formatted info message
	Infof(format string, args ...interface{})
	// Info logs an info message
	Info(args ...interface{})
	// Debugf logs a formatted debug message
	Debugf(format string, args ...interface{})
	// Debug logs a debug message
	Debug(args ...interface{})

	// Create a child logger with fields (this field will be added to the fields
	// of the current logger)
	WithContext(fields map[string]interface{}) Logger
	// Create a child logger with an additional field
	WithField(name, value string) Logger
}

// Level is an enumeration type of the supported loging levels
type Level int

const (
	// Fatal level, is the highest logging level for fatal errors (e.g. panics)
	Fatal Level = iota
	// Error level
	Error
	// Warn level
	Warn
	// Info level
	Info
	// Debug level
	Debug
)

// Options required to configure a Logger implementation
type Options struct {
	Backend   string
	LogFile   string
	LogFormat string
	Level     Level
	Output    io.Writer
}

// GetOutput gets the output where the log is written to
func (o Options) GetOutput() io.Writer {
	if o.Output != nil {
		return o.Output
	}
	if o.LogFile != "" {
		if err := os.MkdirAll(path.Dir(o.LogFile), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create path %s: %s", path.Dir(o.LogFile), err.Error())
		} else {
			file, err := os.OpenFile(o.LogFile, syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open log file %s: %s", o.LogFile, err.Error())
			} else {
				return file
			}
		}
	}
	return os.Stdout
}

// DefaultOptions common default options
var DefaultOptions = &Options{
	Backend: "simple",
}

// NewLogger configure a Logger implementation
func NewLogger(options Options) Logger {
	if factory, ok := backends[options.Backend]; ok {
		return factory(options)
	}
	return NewSimpleLogger(options)
}

type SimpleStackTracer interface {
	ErrorStack() string
}
