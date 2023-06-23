package logging

import (
	"github.com/adharshmk96/stk/utils"
	"github.com/sirupsen/logrus"
)

const (
	EnvLogLevel = "LOG_LEVEL"
)

func setLoggingLevel(logger *logrus.Logger) {
	logLevel := utils.GetEnvOrDefault(EnvLogLevel, "info")

	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

func NewLogrusLogger() *logrus.Logger {
	logger := logrus.New()

	setLoggingLevel(logger)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}
