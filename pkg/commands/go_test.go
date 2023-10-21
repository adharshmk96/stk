package commands_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/pkg/commands"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func TestIsGoMod(t *testing.T) {

	goCmd := commands.NewGoCmd()

	t.Run("returns true if go.mod exists", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)
		defer removeTemp()
		os.Chdir(tempDir)

		goCmd.ModInit("some-package-name")
		assert.True(t, goCmd.IsMod())
		assert.FileExists(t, "go.mod")
	})

	t.Run("returns false if go.mod does not exist", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)
		defer removeTemp()
		os.Chdir(tempDir)

		assert.False(t, goCmd.IsMod())
	})

}

func TestInitGoMod(t *testing.T) {
	goCmd := commands.NewGoCmd()
	t.Run("initializes go module", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)
		defer removeTemp()
		os.Chdir(tempDir)

		err := goCmd.ModInit("some-package-name")
		assert.NoError(t, err)
		assert.FileExists(t, "go.mod")
	})

	t.Run("returns error if go module is already initialized", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)
		defer removeTemp()
		os.Chdir(tempDir)

		err := goCmd.ModInit("some-package-name")
		assert.NoError(t, err)
		assert.FileExists(t, "go.mod")

		err = goCmd.ModInit("some-package-name")
		assert.Error(t, err)
	})

}
