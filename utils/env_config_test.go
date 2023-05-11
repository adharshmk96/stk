package utils_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	t.Run("existing environment variable", func(t *testing.T) {
		os.Setenv("TEST_ENV", "test_value")
		defer os.Unsetenv("TEST_ENV")

		value := utils.GetEnvOrDefault("TEST_ENV", "default_value")
		assert.Equal(t, "test_value", value, "value should be equal to the set environment variable")
	})

	t.Run("non-existing environment variable", func(t *testing.T) {
		value := utils.GetEnvOrDefault("NON_EXISTING_ENV", "default_value")
		assert.Equal(t, "default_value", value, "value should be equal to the default value")
	})
}
