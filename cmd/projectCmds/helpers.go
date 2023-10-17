package projectCmds

import (
	"log"
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
		log.Println("Using repository name as package name.")
		return repoName
	}

	packageName, err := getPackageNameFromGoMod()
	if err == nil && packageName != "" {
		log.Println("Using go module name as package name.")
		return packageName
	}

	argName := getPackageNameFromArg(args)
	if argName != "" {
		log.Println("Using argument as package name.")
		return argName
	}

	configName := viper.GetString("project.package")
	if configName != "" {
		log.Println("Using config project.package as package name.")
		return configName
	}

	randomName := project.RandomName()
	log.Println("Using random name as package name.")
	return randomName
}

func getPackageNameFromArg(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return args[0]
}

func getModuleNameFromArgs(args []string) string {
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

func getPackageNameFromGoMod() (string, error) {
	if !project.IsGoModule() {
		return "", nil
	}

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
