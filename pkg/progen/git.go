package progen

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adharshmk96/stk/pkg/progen/tpl"
)

func initializePackageWithGit(config *Config) error {
	if config.IsGitRepo {
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

func initGit() error {
	return runCommand("git", "init")
}

func IsGitRepo() bool {
	_, err := os.Stat(".git")
	return !os.IsNotExist(err)
}

// getRepoName retrieves the name of the current git repository.
func getRepoName() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	repoUrl := strings.TrimSuffix(string(out), ".git\n")
	repoUrl = strings.ReplaceAll(repoUrl, "https://", "")
	repoUrl = strings.ReplaceAll(repoUrl, "git@", "")
	repoUrl = strings.ReplaceAll(repoUrl, ":", "/")

	return repoUrl, nil
}
