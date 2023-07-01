/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version of stk",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("STK Version:" + version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
