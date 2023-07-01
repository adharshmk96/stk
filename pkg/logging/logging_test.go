package logging_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLoggerLogrus(t *testing.T) {
	t.Run("logger is not nil", func(t *testing.T) {
		logger := logging.NewLogrusLogger()
		assert.NotNil(t, logger, "logger should not be nil")
	})

	t.Run("set log level to 'error'", func(t *testing.T) {
		os.Setenv(logging.ENV_LOG_LEVEL, "error")
		logger := logging.NewLogrusLogger()
		assert.NotNil(t, logger, "logger should not be nil")

		// Get current logger level
		level := logger.IsLevelEnabled(logrus.ErrorLevel)
		assert.True(t, level, "log level should be 'error'")

		infoLevel := logger.IsLevelEnabled(logrus.InfoLevel)
		assert.False(t, infoLevel, "log level should not be 'info'")
	})

	t.Run("set log level to 'debug'", func(t *testing.T) {

		os.Setenv(logging.ENV_LOG_LEVEL, "debug")
		defer os.Unsetenv(logging.ENV_LOG_LEVEL)

		logger := logging.NewLogrusLogger()
		assert.NotNil(t, logger, "logger should not be nil")

		// Get current logger level
		level := logger.IsLevelEnabled(logrus.DebugLevel)
		assert.True(t, level, "log level should be 'debug'")
	})

	t.Run("default log level", func(t *testing.T) {

		os.Setenv(logging.ENV_LOG_LEVEL, "")
		defer os.Unsetenv(logging.ENV_LOG_LEVEL)

		logger := logging.NewLogrusLogger()
		assert.NotNil(t, logger, "logger should not be nil")

		// Get current logger level
		level := logger.IsLevelEnabled(logrus.InfoLevel)
		assert.True(t, level, "default log level should be 'info'")

		leveldebug := logger.IsLevelEnabled(logrus.DebugLevel)
		assert.False(t, leveldebug, "log level should be 'debug'")
	})
}
