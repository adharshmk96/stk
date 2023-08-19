package projectCmds

import (
	"os"
	"os/exec"
	"strings"

	"github.com/iancoleman/strcase"
)

func getPackageNameFromArg(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return args[0]
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

func getAppNameFromPkgName(s string) string {
	lastSlash := strings.LastIndex(s, "/")
	lastPart := s
	if lastSlash != -1 {
		lastPart = s[lastSlash+1:]
	}
	name := strings.ReplaceAll(lastPart, "-", "")
	return strcase.ToLowerCamel(name)
}

func openDirectory(workDir string) error {
	os.MkdirAll(workDir, 0755)

	err := os.Chdir(workDir)
	if err != nil {
		return err
	}
	return nil
}
