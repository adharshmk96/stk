/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package projectCmds

import (
	"log"

	"github.com/adharshmk96/stk/pkg/project"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var deleteModule bool

var ModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "Generate a module for project with gsk and clean architecture.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		isGoModule := project.IsGoModule()
		isGitRepo := project.IsGitRepo()
		pkg := getPackageName(args)
		app := getAppNameFromPkgName(pkg)
		module := getModuleNameFromArgs(args)
		if module == "" {
			log.Fatal("Module name is required.")
			return
		}

		workdir := "."
		err := openDirectory(workdir)
		if err != nil {
			log.Fatal(err)
			return
		}

		modConfig := &project.Config{
			RootPath:     workdir,
			PkgName:      pkg,
			AppName:      app,
			ModName:      strcase.ToLowerCamel(module),
			ExportedName: strcase.ToCamel(module),
			IsGoModule:   isGoModule,
			IsGitRepo:    isGitRepo,
		}

		generator := project.NewGenerator(modConfig)
		if deleteModule {
			log.Println("Deleting module files...")
			err = generator.DeleteModule()
			if err != nil {
				log.Fatal(err)
				return
			}

			log.Println("Deleted module files.")
			return
		}

		log.Println("Generating module files...")
		err = generator.GenerateModule()
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Generated module files.")
	},
}

func init() {
	ModuleCmd.Flags().BoolVarP(&deleteModule, "delete", "d", false, "Delete module files.")
}
