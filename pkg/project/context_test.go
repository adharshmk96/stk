package project_test

import (
	"os"
	"path"
	"testing"

	"github.com/adharshmk96/stk/pkg/commands"
	"github.com/adharshmk96/stk/pkg/project"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func TestWriteDefaultConfig(t *testing.T) {
	t.Run("write default config", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()

		ctx := &project.Context{
			WorkDir: tempDir,
		}

		err := ctx.WriteDefaultConfig()
		assert.NoError(t, err)

		assert.FileExists(t, path.Join(tempDir, project.CONFIG_FILENAME))

	})
}

func TestNewContext(t *testing.T) {
	t.Run("create new context from empty directory", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()

		os.Chdir(tempDir)

		ctx := project.NewContext([]string{"some-package-name"})
		ctx.WorkDir = tempDir

		assert.Equal(t, "some-package-name", ctx.PackageName)
		assert.Equal(t, "somePackageName", ctx.AppName)
		assert.Equal(t, []string{"ping"}, ctx.Modules)
		assert.Equal(t, tempDir, ctx.WorkDir)
		assert.Equal(t, false, ctx.IsGitRepo)
		assert.Equal(t, false, ctx.IsGoModule)
	})

	t.Run("create new context from git repo", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()

		os.Chdir(tempDir)

		gitCmd := commands.NewGitCmd()
		gitCmd.Init()

		ctx := project.NewContext([]string{"some-package-name"})
		ctx.WorkDir = tempDir
		ctx.GitCmd = gitCmd

		assert.Equal(t, "some-package-name", ctx.PackageName)
		assert.Equal(t, "somePackageName", ctx.AppName)
		assert.Equal(t, []string{"ping"}, ctx.Modules)
		assert.Equal(t, tempDir, ctx.WorkDir)
		assert.Equal(t, true, ctx.IsGitRepo)
		assert.Equal(t, false, ctx.IsGoModule)
	})

	t.Run("create new context from go module", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()

		os.Chdir(tempDir)

		goCmd := commands.NewGoCmd()
		err := goCmd.ModInit("github.com/user/package")
		assert.NoError(t, err)

		ctx := project.NewContext([]string{"some-package-name"})
		ctx.WorkDir = tempDir
		ctx.GoCmd = goCmd

		assert.Equal(t, "github.com/user/package", ctx.PackageName)
		assert.Equal(t, "package", ctx.AppName)
		assert.Equal(t, []string{"ping"}, ctx.Modules)
		assert.Equal(t, tempDir, ctx.WorkDir)
		assert.Equal(t, false, ctx.IsGitRepo)
		assert.Equal(t, true, ctx.IsGoModule)
	})

	t.Run("create new context from git repo and go module", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()

		os.Chdir(tempDir)

		gitCmd := commands.NewGitCmd()
		gitCmd.Init()

		goCmd := commands.NewGoCmd()
		err := goCmd.ModInit("github.com/user/package")
		assert.NoError(t, err)

		ctx := project.NewContext([]string{"some-package-name"})
		ctx.WorkDir = tempDir
		ctx.GitCmd = gitCmd
		ctx.GoCmd = goCmd

		assert.Equal(t, "github.com/user/package", ctx.PackageName)
		assert.Equal(t, "package", ctx.AppName)
		assert.Equal(t, []string{"ping"}, ctx.Modules)
		assert.Equal(t, tempDir, ctx.WorkDir)
		assert.Equal(t, true, ctx.IsGitRepo)
		assert.Equal(t, true, ctx.IsGoModule)
	})
}
