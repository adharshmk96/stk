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

// addCmd represents the project command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add project components like modules, boilerplate, etc.",
	Run: func(cmd *cobra.Command, args []string) {

		moduleName := rmModule

		ctx := project.NewContext(args)

		fmt.Println("removing module boilerplate...")
		err := project.GenerateModuleBoilerplate(ctx, moduleName)
		if err != nil {
			fmt.Printf("error while deleting: %s", err)
			return
		}
		fmt.Println("module boilerplate deleted successfully.")
	},
}

func init() {
	addCmd.PersistentFlags().StringVarP(&workDir, "workdir", "w", ".", "project directory")
	addCmd.Flags().StringVarP(&rmModule, "module", "m", "", "module to remove")

	viper.BindPFlag("project.workdir", addCmd.PersistentFlags().Lookup("workdir"))

	rootCmd.AddCommand(addCmd)

}
