package project

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adharshmk96/stk/pkg/project/tpl"
)

func initializePackageWithGit(config *Config) error {
	if _, err := os.Stat(filepath.Join(config.RootPath, ".git")); err == nil {
		log.Println("Directory is already a Git repository.")
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

func getRepoName() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	repoUrl := string(out)
	repoUrl = strings.TrimSuffix(repoUrl, ".git\n")
	repoUrl = strings.ReplaceAll(repoUrl, "https://", "")
	repoUrl = strings.ReplaceAll(repoUrl, "git@", "")
	repoUrl = strings.ReplaceAll(repoUrl, ":", "/")

	return repoUrl, nil
}
