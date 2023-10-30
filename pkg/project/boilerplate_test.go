package project_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/adharshmk96/stk/mocks"
	"github.com/adharshmk96/stk/pkg/commands"
	"github.com/adharshmk96/stk/pkg/project"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func getMockCommands(t *testing.T) (*mocks.GoCmd, *mocks.GitCmd) {
	goCmd := mocks.NewGoCmd(t)
	gitCmd := mocks.NewGitCmd(t)

	return goCmd, gitCmd
}

func TestGenerateBoilerplate(t *testing.T) {

	goCmd := commands.NewGoCmd()
	gitCmd := commands.NewGitCmd()

	t.Run("generates boilerplate", func(t *testing.T) {

		tc := []struct {
			name       string
			isGitRepo  bool
			isGoModule bool
		}{
			{
				name:       "empty project",
				isGitRepo:  false,
				isGoModule: false,
			},
			{
				name:       "git repo",
				isGitRepo:  true,
				isGoModule: false,
			},
			{
				name:       "go module",
				isGitRepo:  false,
				isGoModule: true,
			},
			{
				name:       "git repo and go module",
				isGitRepo:  true,
				isGoModule: true,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				tempDir, removeDir := testutils.SetupTempDirectory(t)
				defer removeDir()

				goCmd, gitCmd := getMockCommands(t)

				if !tt.isGoModule {
					goCmd.On("ModInit", "github.com/sample/sapp").Return(nil)
				}
				if !tt.isGitRepo {
					gitCmd.On("Init").Return(nil)
				}
				goCmd.On("ModTidy").Return(nil)

				ctx := &project.Context{
					PackageName: "github.com/sample/sapp",
					AppName:     "sapp",
					IsGitRepo:   tt.isGitRepo,
					IsGoModule:  tt.isGoModule,
					WorkDir:     tempDir,

					GoCmd:  goCmd,
					GitCmd: gitCmd,
				}

				err := project.GenerateProjectBoilerplate(ctx)
				assert.NoError(t, err)

				assert.FileExists(t, filepath.Join(tempDir, "main.go"))
			})
		}

	})

	t.Run("generates project boilerplate non-mock", func(t *testing.T) {
		tempDir, removeDir := testutils.SetupTempDirectory(t)
		defer removeDir()

		ctx := &project.Context{
			PackageName: "github.com/sample/sapp",
			AppName:     "sapp",
			IsGitRepo:   false,
			IsGoModule:  false,
			WorkDir:     tempDir,

			GoCmd:  goCmd,
			GitCmd: gitCmd,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.NoError(t, err)

		assert.FileExists(t, filepath.Join(tempDir, "main.go"))
		assert.FileExists(t, filepath.Join(tempDir, "go.mod"))
		assert.FileExists(t, filepath.Join(tempDir, "go.sum"))

		assert.True(t, goCmd.IsMod())
		assert.True(t, gitCmd.IsRepo())

	})

}

func TestGenerateModuleBoilerplate(t *testing.T) {
	goCmd := commands.NewGoCmd()
	gitCmd := commands.NewGitCmd()
	t.Run("generates module boilerplate", func(t *testing.T) {
		tempDir, removeDir := testutils.SetupTempDirectory(t)
		defer removeDir()

		goCmd.ModInit("github.com/sample/sapp")
		gitCmd.Init()

		ctx := &project.Context{
			PackageName: "github.com/sample/sapp",
			AppName:     "sapp",
			IsGitRepo:   true,
			IsGoModule:  true,
			WorkDir:     tempDir,

			GoCmd:  goCmd,
			GitCmd: gitCmd,
		}

		err := project.GenerateModuleBoilerplate(ctx, "admin")
		assert.NoError(t, err)

		assert.FileExists(t, filepath.Join(tempDir, "internals/core/entity", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/core/serr", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/service", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/service_test", "admin_test.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/handler", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/helpers", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/transport", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/handler_test", "admin_test.go"))

		assert.True(t, goCmd.IsMod())
		assert.True(t, gitCmd.IsRepo())
	})

	t.Run("errors when go mod init fails", func(t *testing.T) {
		tempDir, removeDir := testutils.SetupTempDirectory(t)
		defer removeDir()

		goCmd, gitCmd := getMockCommands(t)

		gitCmd.On("Init").Return(errors.New(("init failed")))

		ctx := &project.Context{
			PackageName: "github.com/sample/sapp",
			AppName:     "sapp",
			IsGitRepo:   false,
			IsGoModule:  false,
			WorkDir:     tempDir,

			GoCmd:  goCmd,
			GitCmd: gitCmd,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.Error(t, err)
	})

	t.Run("errors when git init fails", func(t *testing.T) {
		tempDir, removeDir := testutils.SetupTempDirectory(t)
		defer removeDir()

		goCmd, gitCmd := getMockCommands(t)

		gitCmd.On("Init").Return(errors.New("init failed"))

		ctx := &project.Context{
			PackageName: "github.com/sample/sapp",
			AppName:     "sapp",
			IsGitRepo:   false,
			IsGoModule:  false,
			WorkDir:     tempDir,

			GoCmd:  goCmd,
			GitCmd: gitCmd,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.Error(t, err)
	})

	t.Run("errors when go mod tidy fails", func(t *testing.T) {
		tempDir, removeDir := testutils.SetupTempDirectory(t)
		defer removeDir()

		goCmd, gitCmd := getMockCommands(t)

		goCmd.On("ModInit", "github.com/sample/sapp").Return(nil)
		gitCmd.On("Init").Return(nil)
		goCmd.On("ModTidy").Return(errors.New("some error"))

		ctx := &project.Context{
			PackageName: "github.com/sample/sapp",
			AppName:     "sapp",
			IsGitRepo:   false,
			IsGoModule:  false,
			WorkDir:     tempDir,

			GoCmd:  goCmd,
			GitCmd: gitCmd,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.Error(t, err)
	})
}

func TestDeleteModuleBoilerplate(t *testing.T) {
	t.Run("deletes module boilerplate", func(t *testing.T) {
		tempDir, removeDir := testutils.SetupTempDirectory(t)
		defer removeDir()

		goCmd := commands.NewGoCmd()
		gitCmd := commands.NewGitCmd()

		goCmd.ModInit("github.com/sample/sapp")
		gitCmd.Init()

		ctx := &project.Context{
			PackageName: "github.com/sample/sapp",
			AppName:     "sapp",
			IsGitRepo:   true,
			IsGoModule:  true,
			WorkDir:     tempDir,

			GoCmd:  goCmd,
			GitCmd: gitCmd,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.NoError(t, err)

		err = project.GenerateModuleBoilerplate(ctx, "admin")
		assert.NoError(t, err)

		assert.FileExists(t, filepath.Join(tempDir, "internals/core/entity", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/core/serr", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/service", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/service_test", "admin_test.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/handler", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/helpers", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/transport", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/handler_test", "admin_test.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/storage/adminStorage", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/storage/adminStorage", "adminQueries.go"))

		assert.DirExists(t, filepath.Join(tempDir, "internals/storage/adminStorage"))

		err = project.DeleteModuleBoilerplate(ctx, "admin")
		assert.NoError(t, err)

		assert.NoFileExists(t, filepath.Join(tempDir, "internals/core/entity", "admin.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/core/serr", "admin.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/service", "admin.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/service_test", "admin_test.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/http/handler", "admin.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/http/helpers", "admin.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/http/transport", "admin.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/http/handler_test", "admin_test.go"))

		assert.NoDirExists(t, filepath.Join(tempDir, "internals/storage/adminStorage"))

		assert.True(t, goCmd.IsMod())
		assert.True(t, gitCmd.IsRepo())
	})

	t.Run("remove default module", func(t *testing.T) {
		tempDir, removeDir := testutils.SetupTempDirectory(t)
		defer removeDir()

		goCmd := commands.NewGoCmd()
		gitCmd := commands.NewGitCmd()

		goCmd.ModInit("github.com/sample/sapp")
		gitCmd.Init()

		ctx := &project.Context{
			PackageName: "github.com/sample/sapp",
			AppName:     "sapp",
			IsGitRepo:   true,
			IsGoModule:  true,
			WorkDir:     tempDir,

			GoCmd:  goCmd,
			GitCmd: gitCmd,
		}

		err := project.GenerateProjectBoilerplate(ctx)
		assert.NoError(t, err)

		err = project.GenerateModuleBoilerplate(ctx, "admin")
		assert.NoError(t, err)

		assert.FileExists(t, filepath.Join(tempDir, "internals/core/entity", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/core/serr", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/service", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/service_test", "admin_test.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/handler", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/helpers", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/transport", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/http/handler_test", "admin_test.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/storage/adminStorage", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/storage/adminStorage", "adminQueries.go"))

		assert.DirExists(t, filepath.Join(tempDir, "internals/storage/adminStorage"))

		err = project.DeleteModuleBoilerplate(ctx, "ping")
		assert.NoError(t, err)

		assert.NoFileExists(t, filepath.Join(tempDir, "internals/core/entity", "ping.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/core/serr", "ping.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/service", "ping.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/service_test", "ping_test.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/http/handler", "ping.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/http/helpers", "ping.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/http/transport", "ping.go"))
		assert.NoFileExists(t, filepath.Join(tempDir, "internals/http/handler_test", "ping_test.go"))

		assert.NoDirExists(t, filepath.Join(tempDir, "internals/storage/pingStorage"))

		assert.True(t, goCmd.IsMod())
		assert.True(t, gitCmd.IsRepo())

	})
}
