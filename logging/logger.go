package logging

import (
	"encoding/json"

	"github.com/adharshmk96/stk/utils"
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
