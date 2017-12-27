package logging

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/pkg/errors"
)

type loggerSimple struct {
	errorLogger *log.Logger
	warnLogger  *log.Logger
	infoLogger  *log.Logger
	debugLogger *log.Logger
	out         io.Writer
}

func init() {
	RegisterBackend("simple", NewSimpleLogger)
	RegisterBackend("null", NewSimpleLoggerNull)
}

// NewSimpleLogger creates a simple Logger based on the golang log package
func NewSimpleLogger(options Options) Logger {
	out := options.GetOutput()
	var errorLogger *log.Logger
	var warnLogger *log.Logger
	var infoLogger *log.Logger
	var debugLogger *log.Logger

	if options.Level >= Error {
		errorLogger = log.New(out, "ERROR: ", log.LstdFlags)
	}
	if options.Level >= Warn {
		warnLogger = log.New(out, "WARN: ", log.LstdFlags)
	}
	if options.Level >= Info {
		infoLogger = log.New(out, "INFO: ", log.LstdFlags)
	}
	if options.Level >= Debug {
		debugLogger = log.New(out, "DEBUG: ", log.LstdFlags)
	}
	return &loggerSimple{
		errorLogger: errorLogger,
		warnLogger:  warnLogger,
		infoLogger:  infoLogger,
		debugLogger: debugLogger,
		out:         out,
	}
}

// NewSimpleLoggerNull create a simple Logger discarding all log entries (i.e.
// /dev/null). Useful for testing where you do not want to polute testing
// output with log messages.
func NewSimpleLoggerNull(options Options) Logger {
	return &loggerSimple{}
}

func (l *loggerSimple) ErrorErr(err error) {
	if l.errorLogger != nil {
		switch richErr := err.(type) {
		case fmt.Formatter:
			l.errorLogger.Printf("%+v", richErr)
		case SimpleStackTracer:
			l.errorLogger.Print(richErr.ErrorStack())
		default:
			wrapped := errors.Wrap(err, err.Error())

			l.errorLogger.Printf("%+v", wrapped)
		}
	}
}

func (l *loggerSimple) Errorf(format string, args ...interface{}) {
	if l.errorLogger != nil {
		l.errorLogger.Printf(format, args...)
	}
}

func (l *loggerSimple) Error(args ...interface{}) {
	if l.errorLogger != nil {
		l.errorLogger.Print(args...)
	}
}

func (l *loggerSimple) Warnf(format string, args ...interface{}) {
	if l.warnLogger != nil {
		l.warnLogger.Printf(format, args...)
	}
}

func (l *loggerSimple) Warn(args ...interface{}) {
	if l.warnLogger != nil {
		l.warnLogger.Print(args...)
	}
}

func (l *loggerSimple) Infof(format string, args ...interface{}) {
	if l.infoLogger != nil {
		l.infoLogger.Printf(format, args...)
	}
}

func (l *loggerSimple) Info(args ...interface{}) {
	if l.infoLogger != nil {
		l.infoLogger.Print(args...)
	}
}

func (l *loggerSimple) Debugf(format string, args ...interface{}) {
	if l.debugLogger != nil {
		l.debugLogger.Printf(format, args...)
	}
}

func (l *loggerSimple) Debug(args ...interface{}) {
	if l.debugLogger != nil {
		l.debugLogger.Print(args...)
	}
}

func (l *loggerSimple) WithContext(fields map[string]interface{}) Logger {
	elements := make([]string, 0)
	for k, v := range fields {
		elements = append(elements, fmt.Sprintf("%s=%v", k, v))
	}

	return l.subLogger(strings.Join(elements, " "))
}

func (l *loggerSimple) WithField(name, value string) Logger {
	return l.subLogger(fmt.Sprintf("%s=%s", name, value))
}

func (l *loggerSimple) subLogger(additionalPrefix string) Logger {
	var errorLogger *log.Logger
	var warnLogger *log.Logger
	var infoLogger *log.Logger
	var debugLogger *log.Logger

	if l.errorLogger != nil {
		errorLogger = log.New(l.out, l.errorLogger.Prefix()+additionalPrefix+" ", log.LstdFlags)
	}
	if l.warnLogger != nil {
		warnLogger = log.New(l.out, l.warnLogger.Prefix()+additionalPrefix+" ", log.LstdFlags)
	}
	if l.infoLogger != nil {
		infoLogger = log.New(l.out, l.infoLogger.Prefix()+additionalPrefix+" ", log.LstdFlags)
	}
	if l.debugLogger != nil {
		debugLogger = log.New(l.out, l.debugLogger.Prefix()+additionalPrefix+" ", log.LstdFlags)
	}
	return &loggerSimple{
		errorLogger: errorLogger,
		warnLogger:  warnLogger,
		infoLogger:  infoLogger,
		debugLogger: debugLogger,
		out:         l.out,
	}
}
