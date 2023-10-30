package project

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/adharshmk96/stk/pkg/commands"
)

func Clean(output string, err error) (string, error) {
	output = strings.ReplaceAll(strings.Split(output, "\n")[0], "'", "")
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}

func OpenDirectory(workDir string) error {
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return err
	}
	return os.Chdir(workDir)
}

func GitRepoName() (string, error) {
	gitCmd := commands.NewGitCmd()
	remoteUrl, err := Clean(gitCmd.GetRemoteOrigin())
	if err != nil {
		return "", err
	}

	repoUrl := strings.TrimSuffix(remoteUrl, ".git")
	repoUrl = strings.ReplaceAll(repoUrl, "https://", "")
	repoUrl = strings.ReplaceAll(repoUrl, "git@", "")
	repoUrl = strings.ReplaceAll(repoUrl, ":", "/")

	return repoUrl, nil
}

func getFirstArg(args []string) string {
	if len(args) > 0 {
		return strings.TrimSpace(args[0])
	}
	return ""
}

func RandomName() string {
	nouns := []string{"apple", "ball", "cat", "dog", "elephant", "fish", "gorilla", "horse", "iguana", "jellyfish", "kangaroo"}
	adjectives := []string{"angry", "big", "cold", "dark", "fast", "good", "happy", "jolly", "kind", "little", "merry", "nice"}

	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)

	return fmt.Sprintf("%s%s", adjectives[randGen.Intn(len(adjectives))], nouns[randGen.Intn(len(nouns))])
}
