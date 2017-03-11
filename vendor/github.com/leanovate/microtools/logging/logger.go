package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"syscall"
)

type Logger interface {
	ErrorErr(error)
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})

	WithContext(fields map[string]interface{}) Logger
	WithField(name, value string) Logger
}

type Level int

const (
	Fatal Level = iota
	Error
	Warn
	Info
	Debug
)

type Options struct {
	Backend   string
	LogFile   string
	LogFormat string
	Level     Level
	Output    io.Writer
}

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

var DefaultOptions = &Options{
	Backend: "logrus",
}

func NewLogger(options Options) Logger {
	switch options.Backend {
	case "simple":
		return NewSimpleLogger(options)
	case "null":
		return NewSimpleLoggerNull()
	}
	return NewLogrusLogger(options)
}

type simpleStackTracer interface {
	ErrorStack() string
}
