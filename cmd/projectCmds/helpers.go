package projectCmds

import (
	"os"
	"os/exec"
	"strings"

	"github.com/adharshmk96/stk/pkg/project"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
)

func getPackageName(args []string) string {
	repoName, err := getRepoName()
	if err != nil && repoName != "" {
		return repoName
	}

	argName := getPackageNameFromArg(args)
	if argName != "" {
		return argName
	}

	configName := viper.GetString("project.package")
	if configName != "" {
		return configName
	}

	randomName := project.RandomName()
	return randomName
}

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

func isGoModule() bool {
	_, err := os.Stat("go.mod")
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getPackageNameFromGoMod() (string, error) {
	cmd := exec.Command("go", "list", "-m")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	pkgName := string(out)
	pkgName = strings.TrimSuffix(pkgName, "\n")

	return pkgName, nil
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
