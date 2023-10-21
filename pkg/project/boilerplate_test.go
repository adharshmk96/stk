package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adharshmk96/stk/pkg/git"
	"github.com/adharshmk96/stk/pkg/project"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateBoilerplate(t *testing.T) {
	t.Run("generates boilerplate in an empty folder", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		os.Chdir(tempDir)

		ctx := &project.Context{
			PackageName: "github.com/adharshmk96/stk",
			AppName:     "stk",
			Modules:     []string{"ping"},
			IsGitRepo:   false,
			IsGoModule:  false,
			WorkDir:     tempDir,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.NoError(t, err)

		assert.FileExists(t, filepath.Join(tempDir, "main.go"))
		assert.FileExists(t, filepath.Join(tempDir, "go.mod"))
		assert.FileExists(t, filepath.Join(tempDir, "go.sum"))

		assert.True(t, project.IsGoModule())
		assert.True(t, git.IsRepo())

	})

	t.Run("generates boilerplate in an existing git repo", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		os.Chdir(tempDir)

		ctx := &project.Context{
			PackageName: "github.com/adharshmk96/stk",
			AppName:     "stk",
			Modules:     []string{"ping"},
			IsGitRepo:   true,
			IsGoModule:  false,
			WorkDir:     tempDir,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.NoError(t, err)

		assert.FileExists(t, filepath.Join(tempDir, "main.go"))
		assert.FileExists(t, filepath.Join(tempDir, "go.mod"))
		assert.FileExists(t, filepath.Join(tempDir, "go.sum"))

		assert.True(t, project.IsGoModule())
		assert.True(t, git.IsRepo())
	})

	t.Run("generates boilerplate in an existing go module", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		os.Chdir(tempDir)

		ctx := &project.Context{
			PackageName: "github.com/adharshmk96/stk",
			AppName:     "stk",
			Modules:     []string{"ping"},
			IsGitRepo:   false,
			IsGoModule:  true,
			WorkDir:     tempDir,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.NoError(t, err)

		assert.FileExists(t, filepath.Join(tempDir, "main.go"))
		assert.FileExists(t, filepath.Join(tempDir, "go.mod"))
		assert.FileExists(t, filepath.Join(tempDir, "go.sum"))

		assert.True(t, project.IsGoModule())
		assert.True(t, git.IsRepo())
	})

	t.Run("generates boilerplate in an existing go module and git repo", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		os.Chdir(tempDir)

		ctx := &project.Context{
			PackageName: "github.com/adharshmk96/stk",
			AppName:     "stk",
			Modules:     []string{"ping"},
			IsGitRepo:   true,
			IsGoModule:  true,
			WorkDir:     tempDir,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.NoError(t, err)

		assert.FileExists(t, filepath.Join(tempDir, "main.go"))
		assert.FileExists(t, filepath.Join(tempDir, "go.mod"))
		assert.FileExists(t, filepath.Join(tempDir, "go.sum"))

		assert.True(t, project.IsGoModule())
		assert.True(t, git.IsRepo())
	})
}
