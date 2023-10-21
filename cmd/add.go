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

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: `add a new module to the project.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := project.NewContext(args)

		fmt.Println("generating module boilerplate...")
		err := project.GenerateModuleBoilerplate(ctx, args[0])
		if err != nil {
			fmt.Printf("error while generating: %s", err)
			return
		}
		fmt.Println("module boilerplate generated successfully.")

	},
}

// addCmd represents the project command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add project components like modules, boilerplate, etc.",
}

func init() {
	addCmd.PersistentFlags().StringVarP(&workDir, "workdir", "w", ".", "project directory")

	viper.BindPFlag("project.workdir", addCmd.PersistentFlags().Lookup("workdir"))

	addCmd.AddCommand(moduleCmd)

	rootCmd.AddCommand(addCmd)

}
