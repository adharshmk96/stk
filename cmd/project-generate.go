/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"log"

	"github.com/adharshmk96/stk/pkg/progen"
	"github.com/spf13/cobra"
)

// init with git after cd
// git init
// git remote add origin <your_repository_url>
// git fetch
// git checkout -t origin/<your_branch_name>

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a project with gsk and clean architecture.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Generating project files...")

		isGoModule := progen.IsGoModule()
		isGitRepo := progen.IsGitRepo()
		pkg := progen.GetPackageName(args)
		app := progen.GetAppNameFromPkgName(pkg)

		workdir := "."
		err := progen.OpenDirectory(workdir)
		if err != nil {
			log.Fatal(err)
			return
		}

		config := &progen.Config{
			RootPath:     workdir,
			PkgName:      pkg,
			AppName:      app,
			ModName:      "ping",
			ExportedName: "Ping",
			IsGoModule:   isGoModule,
			IsGitRepo:    isGitRepo,
		}

		generator := progen.NewGenerator(config)
		err = generator.GenerateProject()
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Project generated successfully.")
	},
}

func init() {
	projectCmd.AddCommand(generateCmd)
}
