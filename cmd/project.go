/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"github.com/adharshmk96/stk/cmd/project"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectPackageName string

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project structure, create new project, add sections, etc.",
}

func init() {
	projectCmd.PersistentFlags().StringVar(&projectPackageName, "pkg", "", "project package name or repository name")

	viper.BindPFlag("project.package", projectCmd.PersistentFlags().Lookup("pkg"))

	projectCmd.AddCommand(project.ModuleCmd)
	projectCmd.AddCommand(project.GenerateCmd)

	rootCmd.AddCommand(projectCmd)

}
