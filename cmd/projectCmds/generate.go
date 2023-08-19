/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package projectCmds

import (
	"log"

	"github.com/adharshmk96/stk/pkg/project"
	"github.com/spf13/cobra"
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

// init with git after cd
// git init
// git remote add origin <your_repository_url>
// git fetch
// git checkout -t origin/<your_branch_name>

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a project with gsk and clean architecture.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Generating project files...")

		pkg := getPackageName(args)
		app := getAppNameFromPkgName(pkg)

		workdir := "."
		err := openDirectory(workdir)
		if err != nil {
			log.Fatal(err)
			return
		}

		config := &project.Config{
			RootPath: workdir,
			PkgName:  pkg,
			AppName:  app,
		}

		err = project.GenerateProject(config)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Project generated successfully.")
	},
}
