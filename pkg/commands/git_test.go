package commands_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/pkg/commands"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

// type GitCmd interface {
// 	RunCmd(args ...string) (string, error)
// 	Init() error
// 	AddRemote(remoteName, remoteUrl string) error
// 	GetRemoteOrigin() (string, error)
// 	Revparse(ref string) (string, error)

// 	IsRepo() bool
// }

func TestGitCmd(t *testing.T) {
	t.Run("run cmd with args", func(t *testing.T) {
		testDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()
		os.Chdir(testDir)

		gitCmd := commands.NewGitCmd()

		isRepo := gitCmd.IsRepo()
		assert.False(t, isRepo)

		_, err := gitCmd.RunCmd("init")
		assert.NoError(t, err)

		isRepo = gitCmd.IsRepo()
		assert.True(t, isRepo)
		assert.DirExists(t, testDir+"/.git")

	})

	t.Run("remote origin", func(t *testing.T) {
		testDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()
		os.Chdir(testDir)

		gitCmd := commands.NewGitCmd()

		remote, err := gitCmd.GetRemoteOrigin()
		assert.Error(t, err)
		assert.Empty(t, remote)

		_, err = gitCmd.RunCmd("init")
		assert.NoError(t, err)

		err = gitCmd.AddRemote("origin", "https://github.com/test/test.git")
		assert.NoError(t, err)

		remote, err = gitCmd.GetRemoteOrigin()
		assert.NoError(t, err)
		assert.Equal(t, "https://github.com/test/test.git\n", remote)

	})

}
