package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display the version of singlemod",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("singlemod version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
