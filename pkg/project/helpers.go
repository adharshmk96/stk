package project

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/adharshmk96/stk/pkg/git"
)

func RunCmd(cmd string, args ...string) (string, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	c := exec.Command(cmd, args...)
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	return stdout.String(), nil
}

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
	remoteUrl, err := Clean(git.GetRemoteOrigin())
	if err != nil {
		return "", err
	}

	repoUrl := strings.TrimSuffix(remoteUrl, ".git\n")
	repoUrl = strings.ReplaceAll(repoUrl, "https://", "")
	repoUrl = strings.ReplaceAll(repoUrl, "git@", "")
	repoUrl = strings.ReplaceAll(repoUrl, ":", "/")

	return repoUrl, nil
}

func getFirstArg(args []string) string {
	if len(args) > 0 {
		return args[0]
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
