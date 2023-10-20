package progen

import (
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
)

// GetPackageName determines the package name by checking various sources in order of priority.
func GetPackageName(args []string) string {
	if repoName, err := getRepoName(); err == nil {
		log.Println("Using repository name as package name.")
		return repoName
	}

	if packageName, err := getPackageNameFromGoMod(); err == nil {
		log.Println("Using go module name as package name.")
		return packageName
	}

	if argName := getFirstArg(args); argName != "" {
		log.Println("Using argument as package name.")
		return argName
	}

	if configName := viper.GetString("project.package"); configName != "" {
		log.Println("Using config project.package as package name.")
		return configName
	}

	log.Println("Using random name as package name.")
	return RandomName()
}

// getFirstArg returns the first argument from the args slice.
func getFirstArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// getPackageNameFromGoMod retrieves the go module name if present.
func getPackageNameFromGoMod() (string, error) {
	if !IsGoModule() {
		return "", errors.New("not a Go module")
	}

	cmd := exec.Command("go", "list", "-m")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(out), "\n"), nil
}

// getAppNameFromPkgName transforms package name into an app name in lowerCamelCase.
func GetAppNameFromPkgName(s string) string {
	lastPart := s[strings.LastIndex(s, "/")+1:]
	name := strings.ReplaceAll(lastPart, "-", "")
	return strcase.ToLowerCamel(name)
}
