package project

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adharshmk96/stk/pkg/project/tpl"
)

func initializePackageWithGit(config *Config) error {
	if IsGitRepo() {
		return nil
	}

	exec.Command("git", "init", ".").Run()

	gitIgnorePath := filepath.Join(".gitignore")
	file, err := os.Create(gitIgnorePath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write([]byte(tpl.GITIGNORE_TPL.Content))

	return nil
}

func IsGitRepo() bool {
	_, err := os.Stat(".git")
	return !os.IsNotExist(err)
}

func IsGoModule() bool {
	_, err := os.Stat("go.mod")
	return !os.IsNotExist(err)
}
