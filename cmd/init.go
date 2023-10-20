/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/adharshmk96/stk/pkg/project"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func writeDefaultConfig(ctx *project.Context) {
	fmt.Println("initializing config file...")

	// project configs
	viper.Set("name", ctx.AppName)
	viper.Set("version", "v0.0.1")
	viper.Set("description", "This project is generated using stk.")
	viper.Set("author", "")

	// module configs
	viper.Set("project.modules", ctx.Modules)

	// Migrator configs
	viper.Set("migrator.workdir", "./stk-migrations")
	viper.Set("migrator.database", "sqlite3")
	viper.Set("migrator.sqlite.filepath", "stk.db")

	// Create the config file
	err := viper.WriteConfigAs(".stk.yaml")
	if err != nil {
		fmt.Printf("error while writing config file: %s", err)
	}

	fmt.Println("default configs written successfully.")

}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a new stk project",
	Long:  `a new project will be created in the current directory initializing the go module, git, boilerplate code required to start a new project.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := project.NewContext(args)
		writeDefaultConfig(ctx)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
