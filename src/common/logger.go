/*
Logger implementation, using this as is from
github.com/hpcloud/cf-usb/cmd/usb/logger.go
*/

package common

import (
	"os"

	"github.com/pivotal-golang/lager"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

func NewLogger(level string) lager.Logger {
	var logger = lager.NewLogger("csm")

	var minLogLevel lager.LogLevel
	switch level {
	case DEBUG:
		minLogLevel = lager.DEBUG
	case INFO:
		minLogLevel = lager.INFO
	case ERROR:
		minLogLevel = lager.ERROR
	case FATAL:
		minLogLevel = lager.FATAL
	default:
		minLogLevel = lager.INFO
	}

	logger.RegisterSink(lager.NewWriterSink(os.Stdout, minLogLevel))

	logger.Info("Log level set to:", lager.Data{"level": minLogLevel})

	return logger
}
