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

func getModuleNameFromArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return args[0]
}

var ModuleCmd = &cobra.Command{
	Use:   "module",
	Short: "Generate a module for project with gsk and clean architecture.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Generating module files...")

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

		config := &project.Config{
			RootPath: workdir,
			PkgName:  pkg,
			AppName:  app,
		}

		modConfig := &project.ModuleConfig{
			ModName:      module,
			ExportedName: strcase.ToCamel(module),
		}

		err = project.GenerateModule(config, modConfig)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Generated module files.")
	},
}
