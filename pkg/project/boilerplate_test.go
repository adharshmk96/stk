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

		assert.DirExists(t, filepath.Join(tempDir, "internals/admin"))

		assert.FileExists(t, filepath.Join(tempDir, "internals/admin/domain", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/admin/serr", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/admin/service", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/admin/service", "admin_test.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/admin/api/handler", "admin.go"))
		assert.FileExists(t, filepath.Join(tempDir, "internals/admin/api/transport", "admin.go"))

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

func assertBoilerplateExists(t *testing.T, tempDir, moduleName string) {
	dirs := []string{
		filepath.Join("internals", moduleName),
		filepath.Join("internals", moduleName, "domain"),
		filepath.Join("internals", moduleName, "serr"),
		filepath.Join("internals", moduleName, "service"),
		filepath.Join("internals", moduleName, "web"),
		filepath.Join("internals", moduleName, "storage"),
		filepath.Join("internals", moduleName, "storage"),
		filepath.Join("internals", moduleName, "api/handler"),
		filepath.Join("internals", moduleName, "api/transport"),
	}
	files := []string{
		filepath.Join("internals", moduleName, "routes.go"),
		filepath.Join("internals", moduleName, "service", moduleName+".go"),
		filepath.Join("internals", moduleName, "service", moduleName+"_test.go"),
		filepath.Join("internals", moduleName, "storage", moduleName+".go"),
		filepath.Join("internals", moduleName, "storage", moduleName+"Queries.go"),
		filepath.Join("internals", moduleName, "domain", moduleName+".go"),
		filepath.Join("internals", moduleName, "serr", moduleName+".go"),
		filepath.Join("internals", moduleName, "api/handler", moduleName+".go"),
		filepath.Join("internals", moduleName, "api/handler", moduleName+"_test.go"),
		filepath.Join("internals", moduleName, "api/transport", moduleName+".go"),
		filepath.Join("internals", moduleName, "web", moduleName+".go"),
		filepath.Join("server/routing", moduleName+".go"),
	}

	for _, dir := range dirs {
		assert.DirExists(t, filepath.Join(tempDir, dir))
	}

	for _, file := range files {
		assert.FileExists(t, filepath.Join(tempDir, file))
	}

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

		assertBoilerplateExists(t, tempDir, "admin")

		err = project.DeleteModuleBoilerplate(ctx, "admin")
		assert.NoError(t, err)

		assert.NoDirExists(t, filepath.Join(tempDir, "internals/admin"))
		assert.NoFileExists(t, filepath.Join(tempDir, "server/routing", "admin.go"))

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
		assertBoilerplateExists(t, tempDir, "ping")

		err = project.GenerateModuleBoilerplate(ctx, "admin")
		assert.NoError(t, err)
		assertBoilerplateExists(t, tempDir, "admin")

		err = project.DeleteModuleBoilerplate(ctx, "ping")
		assert.NoError(t, err)

		assert.NoDirExists(t, filepath.Join(tempDir, "internals/ping"))
		assert.NoFileExists(t, filepath.Join(tempDir, "server/routing", "ping.go"))

		assert.True(t, goCmd.IsMod())
		assert.True(t, gitCmd.IsRepo())

	})
}
