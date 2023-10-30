package utils_test

import (
	"testing"

	"github.com/adharshmk96/stk/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetFirst(t *testing.T) {
	t.Run("should return first element", func(t *testing.T) {
		t.Parallel()
		// given
		input := []string{"a", "b", "c"}
		// when
		result := utils.GetFirst(input...)
		// then
		assert.Equal(t, "a", result)
	})
	t.Run("should return empty string if input is empty", func(t *testing.T) {
		t.Parallel()
		// given
		input := []string{}
		// when
		result := utils.GetFirst(input...)
		// then
		assert.Equal(t, "", result)
	})

	t.Run("should return first non-empty string", func(t *testing.T) {
		t.Parallel()
		// given
		input := []string{"", "b", "c"}
		// when
		result := utils.GetFirst(input...)
		// then
		assert.Equal(t, "b", result)
	})
}
