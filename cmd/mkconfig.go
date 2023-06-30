/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// mkconfigCmd represents the mkconfig command
var mkconfigCmd = &cobra.Command{
	Use:   "mkconfig",
	Short: "Generate a config file for stk",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Create("stk.yaml")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
	},
}

func init() {
	rootCmd.AddCommand(mkconfigCmd)
}
