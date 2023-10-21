package project

import (
	"log"
	"strings"

	"github.com/adharshmk96/stk/pkg/commands"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
)

// getAppNameFromPkgName transforms package name into an app name in lowerCamelCase.
func GetAppNameFromPkgName(s string) string {
	name := s[strings.LastIndex(s, "/")+1:]
	name = strings.ReplaceAll(name, "-", "_")
	return strcase.ToLowerCamel(name)
}

// GetPackageName determines the package name by checking various sources in order of priority.
func GetPackageName(args []string) string {
	if repoName, err := GitRepoName(); err == nil {
		log.Println("using repository name as package name.")
		return repoName
	}

	goCmd := commands.NewGoCmd()
	packageName, err := goCmd.ModPackageName()
	if err == nil {
		log.Println("using go module name as package name.")
		return packageName
	}

	if argName := getFirstArg(args); argName != "" {
		log.Println("using argument as package name.")
		return argName
	}

	if configName := viper.GetString("project.package"); configName != "" {
		log.Println("using config project.package as package name.")
		return configName
	}

	log.Println("using random name as package name.")
	return RandomName()
}
