package project_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/pkg/git"
	"github.com/adharshmk96/stk/pkg/project"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetPackageName(t *testing.T) {
	t.Run("gets package name from git repo", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		os.Chdir(tempDir)

		git.Init()
		git.AddRemote("origin", "https://github.com/adharshmk96/stk")

		project.InitGoMod("some-package-name")

		packageName := project.GetPackageName([]string{"just some args"})

		assert.Equal(t, "github.com/adharshmk96/stk", packageName)
	})

	t.Run("gets package name from go.mod", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		os.Chdir(tempDir)

		git.Init()

		project.InitGoMod("some-package-name")

		packageName := project.GetPackageName([]string{"just some args"})
		assert.Equal(t, "some-package-name", packageName)
	})

	t.Run("gets package name from first arg", func(t *testing.T) {
		packageName := project.GetPackageName([]string{"some-package-name"})
		assert.Equal(t, "some-package-name", packageName)
	})

	t.Run("assign random name", func(t *testing.T) {
		packageName := project.GetPackageName([]string{})
		assert.NotEmpty(t, packageName)
	})
}

func TestGetAppNameFromPkgName(t *testing.T) {
	tc := []struct {
		pkgName string
		appName string
	}{
		{"stk", "stk"},
		{"github.com/adharshmk96/stk", "stk"},
		{"github.com/adharshmk96/stk-cli", "stkCli"},
		{"github.com/adharshmk96/stk-cli-go", "stkCliGo"},
		{"github.com/adharshmk96/stk_cli-go", "stkCliGo"},
	}

	for _, c := range tc {
		t.Run(c.pkgName, func(t *testing.T) {
			appName := project.GetAppNameFromPkgName(c.pkgName)
			assert.Equal(t, c.appName, appName)
		})
	}
}
