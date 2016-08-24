/*
Logger implementation, using this as is from
github.com/hpcloud/cf-usb/cmd/usb/logger.go
*/

package common

import (
	"log/syslog"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
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
	sysLogLevel := syslog.LOG_DEBUG
	switch level {
	case INFO:
		minLogLevel = logrus.InfoLevel
		sysLogLevel = syslog.LOG_INFO
	case WARN:
		minLogLevel = logrus.WarnLevel
		sysLogLevel = syslog.LOG_WARNING
	case ERROR:
		minLogLevel = logrus.ErrorLevel
		sysLogLevel = syslog.LOG_ERR
	case FATAL:
		minLogLevel = logrus.FatalLevel
		sysLogLevel = syslog.LOG_CRIT
	case DEBUG:
		minLogLevel = logrus.DebugLevel
		sysLogLevel = syslog.LOG_DEBUG
	}

	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = minLogLevel

	config := NewServiceManagerConfiguration()

	if config.FlightRecorderEndpoint() != ":" {
		hook, err := logrus_syslog.NewSyslogHook("tcp", config.FlightRecorderEndpoint(), sysLogLevel, component)
		if err != nil {
			logger.Warnf("Unable to connect to flight recorder %+v", err)
		} else {
			logger.Hooks.Add(hook)
		}
	} else {
		logger.Info("Flight recorder endpoint not set.")
	}
	return logger
}
