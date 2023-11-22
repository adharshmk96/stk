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

var workDir string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a new stk project",
	Long:  `a new project will be created in the current directory initializing the go module, git, boilerplate code required to start a new project.`,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := cmd.Flag("cfg").Value.String() == "true"

		ctx := project.NewContext(args)

		fmt.Println("initializing config file...")
		err := ctx.WriteDefaultConfig()
		if err != nil {
			fmt.Printf("error while writing config file: %s", err)
		}
		fmt.Println("default configs written successfully.")

		if cfg {
			return
		}

		fmt.Println("generating boilerplate...")
		err = project.GenerateProjectBoilerplate(ctx)
		if err != nil {
			fmt.Printf("error while generating: %s", err)
			return
		}

		fmt.Printf("boilerplate generated successfully. \n\nrun 'make init' to initialize the project with necessary tools.\n\n")
	},
}

func init() {
	initCmd.PersistentFlags().StringVarP(&workDir, "workdir", "w", ".", "project directory")
	initCmd.Flags().Bool("cfg", false, "config file only")

	viper.BindPFlag("project.workdir", initCmd.PersistentFlags().Lookup("workdir"))

	rootCmd.AddCommand(initCmd)
}
