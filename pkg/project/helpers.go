package project

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

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
