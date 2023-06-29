/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package migration

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrationName string

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate migration files",
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("path").Value.String()
		fmt.Println(path)
	},
}

func init() {
	GenerateCmd.Flags().StringVarP(&migrationName, "name", "n", "", "migration name")
}
