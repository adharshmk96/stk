/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/adharshmk96/stk/cmd/migration"
	"github.com/spf13/cobra"
)

var rootFolder string
var database string

// migrationCmd represents the generate command
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "database migration commands",
}

func init() {
	migrationCmd.PersistentFlags().StringVarP(&rootFolder, "path", "p", "./migration", "migration folder path (default ./migration))")
	migrationCmd.PersistentFlags().StringVarP(&database, "database", "d", "sqlite", "database type ( default sqlite )")

	migrationCmd.AddCommand(migration.GenerateCmd)

	rootCmd.AddCommand(migrationCmd)
}
