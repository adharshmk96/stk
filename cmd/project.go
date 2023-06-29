/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project structure, create new project, add sections, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("project command not supported yet")
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
