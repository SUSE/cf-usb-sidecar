/*
Logger implementation, using this as is from
github.com/SUSE/cf-usb/cmd/usb/logger.go
*/

package common

import (
	"os"

	"github.com/Sirupsen/logrus"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
	WARN  = "warn"
)

func NewLogger(level string, component string) *logrus.Logger {
	minLogLevel := logrus.DebugLevel
	switch level {
	case INFO:
		minLogLevel = logrus.InfoLevel
	case WARN:
		minLogLevel = logrus.WarnLevel
	case ERROR:
		minLogLevel = logrus.ErrorLevel
	case FATAL:
		minLogLevel = logrus.FatalLevel
	case DEBUG:
		minLogLevel = logrus.DebugLevel
	}

	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = minLogLevel

	return logger
}
