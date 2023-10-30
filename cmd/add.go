/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/adharshmk96/stk/pkg/project"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the project command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add project components like modules, boilerplate, etc.",
	Run: func(cmd *cobra.Command, args []string) {

		moduleName := strings.TrimSpace(cmd.Flag("module").Value.String())

		if moduleName == "" {
			fmt.Println("module name is required.")
			return
		}

		ctx := project.NewContext(args)

		fmt.Println("removing module boilerplate...")
		err := project.GenerateModuleBoilerplate(ctx, moduleName)
		if err != nil {
			fmt.Printf("error while deleting: %s", err)
			return
		}
		fmt.Println("module boilerplate added successfully.")
	},
}

func init() {
	addCmd.PersistentFlags().StringVarP(&workDir, "workdir", "w", ".", "project directory")
	addCmd.Flags().StringP("module", "m", "", "module name to add")

	viper.BindPFlag("project.workdir", addCmd.PersistentFlags().Lookup("workdir"))

	rootCmd.AddCommand(addCmd)

}
