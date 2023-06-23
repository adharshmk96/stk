package logging_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetLogger(t *testing.T) {
	t.Run("set log level to 'error'", func(t *testing.T) {
		os.Setenv(logging.EnvLogLevel, "error")
		logger := logging.NewZapLogger()
		assert.NotNil(t, logger, "logger should not be nil")

		// Get current logger level
		level := logger.Core().Enabled(zap.ErrorLevel)
		assert.True(t, level, "log level should be 'error'")

		infoLevel := logger.Core().Enabled(zap.InfoLevel)
		assert.False(t, infoLevel, "log level should not be 'info'")
	})

	t.Run("set log level to 'debug'", func(t *testing.T) {

		os.Setenv(logging.EnvLogLevel, "debug")
		defer os.Unsetenv(logging.EnvLogLevel)

		logger := logging.NewZapLogger()
		assert.NotNil(t, logger, "logger should not be nil")

		// Get current logger level
		level := logger.Core().Enabled(zap.DebugLevel)
		assert.True(t, level, "log level should be 'debug'")
	})

	t.Run("default log level", func(t *testing.T) {

		os.Setenv(logging.EnvLogLevel, "")
		defer os.Unsetenv(logging.EnvLogLevel)

		logger := logging.NewZapLogger()
		assert.NotNil(t, logger, "logger should not be nil")

		// Get current logger level
		level := logger.Core().Enabled(zap.InfoLevel)
		assert.True(t, level, "default log level should be 'info'")

		leveldebug := logger.Core().Enabled(zap.DebugLevel)
		assert.False(t, leveldebug, "log level should be 'debug'")
	})

}

func TestGetLoggerLogrus(t *testing.T) {
	t.Run("logger is not nil", func(t *testing.T) {
		logger := logging.NewLogrusLogger()
		assert.NotNil(t, logger, "logger should not be nil")
	})

	t.Run("set log level to 'error'", func(t *testing.T) {
		os.Setenv(logging.EnvLogLevel, "error")
		logger := logging.NewLogrusLogger()
		assert.NotNil(t, logger, "logger should not be nil")

		// Get current logger level
		level := logger.IsLevelEnabled(logrus.ErrorLevel)
		assert.True(t, level, "log level should be 'error'")

		infoLevel := logger.IsLevelEnabled(logrus.InfoLevel)
		assert.False(t, infoLevel, "log level should not be 'info'")
	})

	t.Run("set log level to 'debug'", func(t *testing.T) {

		os.Setenv(logging.EnvLogLevel, "debug")
		defer os.Unsetenv(logging.EnvLogLevel)

		logger := logging.NewLogrusLogger()
		assert.NotNil(t, logger, "logger should not be nil")

		// Get current logger level
		level := logger.IsLevelEnabled(logrus.DebugLevel)
		assert.True(t, level, "log level should be 'debug'")
	})

	t.Run("default log level", func(t *testing.T) {

		os.Setenv(logging.EnvLogLevel, "")
		defer os.Unsetenv(logging.EnvLogLevel)

		logger := logging.NewLogrusLogger()
		assert.NotNil(t, logger, "logger should not be nil")

		// Get current logger level
		level := logger.IsLevelEnabled(logrus.InfoLevel)
		assert.True(t, level, "default log level should be 'info'")

		leveldebug := logger.IsLevelEnabled(logrus.DebugLevel)
		assert.False(t, leveldebug, "log level should be 'debug'")
	})
}
