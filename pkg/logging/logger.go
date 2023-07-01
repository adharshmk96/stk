package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	ENV_LOG_LEVEL = "LOG_LEVEL"

	DEFAULT_LOG_LEVEL = "info"
)

func setLoggingLevel(logger *logrus.Logger) {
	viper.SetDefault(ENV_LOG_LEVEL, DEFAULT_LOG_LEVEL)
	viper.AutomaticEnv()
	logLevel := viper.GetString(ENV_LOG_LEVEL)

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
