package project_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/pkg/git"
	"github.com/adharshmk96/stk/pkg/project"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGitRepoName(t *testing.T) {
	t.Run("return https git repo name", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		git.Init()
		git.AddRemote("origin", "https://github.com/adharshmk96/stk")

		repoName, err := project.GitRepoName()
		assert.NoError(t, err)
		assert.Equal(t, "github.com/adharshmk96/stk", repoName)
	})

	t.Run("return ssh git repo name", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		git.Init()
		git.AddRemote("origin", "git@github.com:adharshmk96/stk")

		repoName, err := project.GitRepoName()
		assert.NoError(t, err)
		assert.Equal(t, "github.com/adharshmk96/stk", repoName)

	})

	t.Run("return error if git repo is not initialized", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		_, err := project.GitRepoName()
		assert.Error(t, err)
	})

	t.Run("return error if git remote is not added", func(t *testing.T) {
		tempDir, removeTemp := testutils.CreateTempDirectory(t)

		defer removeTemp()

		os.Chdir(tempDir)

		git.Init()

		_, err := project.GitRepoName()
		assert.Error(t, err)
	})
}
