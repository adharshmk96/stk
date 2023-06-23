package logging

import (
	"encoding/json"

	"github.com/adharshmk96/stk/utils"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

const (
	EnvLogLevel = "LOG_LEVEL"
)

func NewZapLogger() *zap.Logger {

	var err error
	var logger *zap.Logger

	rawJSON := []byte(`{
		"level": "` + utils.GetEnvOrDefault(EnvLogLevel, "info") + `",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase",
		  "timeKey": "ts",
		  "timeEncoder": "iso8601"
		}
	  }`)

	var cfg zap.Config

	if err = json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	logger, err = cfg.Build()

	if err != nil {
		panic(err)
	}
	return logger
}

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
