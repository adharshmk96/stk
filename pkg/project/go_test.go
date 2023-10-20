package project_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/pkg/project"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func TestIsGoMod(t *testing.T) {
	t.Run("return true if a path is go module", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		project.InitGoMod("new-project")

		goMod := project.IsGoModule()

		assert.True(t, goMod)
	})

	t.Run("return false if a path is not go module", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		goMod := project.IsGoModule()

		assert.False(t, goMod)
	})
}

func TestInitGoMod(t *testing.T) {
	t.Run("initialize go module", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		err := project.InitGoMod("new-project")

		assert.Nil(t, err)
		assert.FileExists(t, "go.mod")
	})

	t.Run("return error if go module is already initialized", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		project.InitGoMod("new-project")

		err := project.InitGoMod("new-project")

		assert.NotNil(t, err)
	})
}

func TestGoModTidy(t *testing.T) {
	t.Run("tidy go module", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		project.InitGoMod("new-project")

		err := project.GoModTidy()

		assert.Nil(t, err)
	})

	t.Run("return error if go module is not initialized", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		err := project.GoModTidy()

		assert.NotNil(t, err)
	})
}

func TestGoModPackageName(t *testing.T) {
	t.Run("return go module package name", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		project.InitGoMod("new-project")

		packageName, err := project.GoModPackageName()

		assert.NoError(t, err)
		assert.Equal(t, "new-project", packageName)
	})

	t.Run("return error if go module is not initialized", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		_packageName, err := project.GoModPackageName()

		assert.Equal(t, "", _packageName)
		assert.Error(t, err)
	})
}
