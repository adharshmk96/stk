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
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove project components like modules, boilerplate, etc.",
	Run: func(cmd *cobra.Command, args []string) {

		moduleName := strings.TrimSpace(cmd.Flag("module").Value.String())

		if moduleName == "" {
			fmt.Println("module name is required.")
			return
		}

		ctx := project.NewContext(args)

		err := project.DeleteModuleBoilerplate(ctx, moduleName)
		if err != nil {
			fmt.Printf("error while deleting: %s", err)
			return
		}
		fmt.Println("module boilerplate deleted successfully.")
	},
}

func init() {
	rmCmd.PersistentFlags().StringVarP(&workDir, "workdir", "w", ".", "project directory")
	rmCmd.Flags().StringP("module", "m", "", "module name to remove")

	viper.BindPFlag("project.workdir", rmCmd.PersistentFlags().Lookup("workdir"))

	rootCmd.AddCommand(rmCmd)

}
