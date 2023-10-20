/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package project

import (
	"log"

	"github.com/adharshmk96/stk/pkg/progen"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var deleteModule bool

func getModuleNameFromArgs(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

var ModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "Generate a module for project with gsk and clean architecture.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		isGoModule := progen.IsGoModule()
		isGitRepo := progen.IsGitRepo()
		pkg := progen.GetPackageName(args)
		app := progen.GetAppNameFromPkgName(pkg)
		module := getModuleNameFromArgs(args)
		if module == "" {
			log.Fatal("Module name is required.")
			return
		}

		workdir := "."
		err := progen.OpenDirectory(workdir)
		if err != nil {
			log.Fatal(err)
			return
		}

		modConfig := &progen.Config{
			RootPath:     workdir,
			PkgName:      pkg,
			AppName:      app,
			ModName:      strcase.ToLowerCamel(module),
			ExportedName: strcase.ToCamel(module),
			IsGoModule:   isGoModule,
			IsGitRepo:    isGitRepo,
		}

		generator := progen.NewGenerator(modConfig)
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
