package logging

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pkg/errors"
)

type loggerSimple struct {
	logger *log.Logger
	level  Level
	out    io.Writer
}

// NewSimpleLogger creates a simple Logger based on the golang log package
func NewSimpleLogger(options Options) Logger {
	out := options.GetOutput()
	return &loggerSimple{
		logger: log.New(out, "", log.LstdFlags),
		level:  options.Level,
		out:    out,
	}
}

// NewSimpleLoggerNull create a simple Logger discarding all log entries (i.e.
// /dev/null). Useful for testing where you do not want to polute testing
// output with log messages.
func NewSimpleLoggerNull() Logger {
	return &loggerSimple{
		logger: log.New(ioutil.Discard, "", log.LstdFlags),
		level:  Fatal,
		out:    ioutil.Discard,
	}
}

func (l *loggerSimple) ErrorErr(err error) {
	if l.level >= Error {
		switch richErr := err.(type) {
		case fmt.Formatter:
			l.logger.Printf("%+v", richErr)
		case simpleStackTracer:
			l.logger.Print(richErr.ErrorStack())
		default:
			wrapped := errors.Wrap(err, err.Error())

			l.logger.Printf("%+v", wrapped)
		}
	}
}

func (l *loggerSimple) Errorf(format string, args ...interface{}) {
	if l.level >= Error {
		l.logger.Printf(format, args...)
	}
}

func (l *loggerSimple) Error(args ...interface{}) {
	if l.level >= Error {
		l.logger.Print(args...)
	}
}

func (l *loggerSimple) Warnf(format string, args ...interface{}) {
	if l.level >= Warn {
		l.logger.Printf(format, args...)
	}
}

func (l *loggerSimple) Warn(args ...interface{}) {
	if l.level >= Warn {
		l.logger.Print(args...)
	}
}

func (l *loggerSimple) Infof(format string, args ...interface{}) {
	if l.level >= Info {
		l.logger.Printf(format, args...)
	}
}

func (l *loggerSimple) Info(args ...interface{}) {
	if l.level >= Info {
		l.logger.Print(args...)
	}
}

func (l *loggerSimple) Debugf(format string, args ...interface{}) {
	if l.level >= Debug {
		l.logger.Printf(format, args...)
	}
}

func (l *loggerSimple) Debug(args ...interface{}) {
	if l.level >= Debug {
		l.logger.Print(args...)
	}
}

func (l *loggerSimple) WithContext(fields map[string]interface{}) Logger {
	elements := make([]string, 0)
	elements = append(elements, l.logger.Prefix())
	for k, v := range fields {
		elements = append(elements, fmt.Sprintf("%s=%v", k, v))
	}

	prefix := strings.Join(elements, " ")
	return &loggerSimple{
		logger: log.New(l.out, prefix, log.LstdFlags),
		level:  l.level,
		out:    l.out,
	}
}

func (l *loggerSimple) WithField(name, value string) Logger {
	prefix := fmt.Sprintf("%s %s=%s", l.logger.Prefix(), name, value)
	return &loggerSimple{
		logger: log.New(l.out, prefix, log.LstdFlags),
		level:  l.level,
		out:    l.out,
	}
}
