/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// verifyCmd represents the project command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify the project structure, check if the project follows required pattern.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("project command not supported yet")
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
