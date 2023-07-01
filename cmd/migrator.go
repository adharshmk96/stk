/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"github.com/adharshmk96/stk/cmd/migCommands"
	"github.com/spf13/cobra"
)

var rootFolder string
var database string

// migratorCmd represents the generate command
var migratorCmd = &cobra.Command{
	Use:   "migrator",
	Short: "database migrator commands",
}

func init() {
	migratorCmd.PersistentFlags().StringVarP(&rootFolder, "path", "p", "./stk-migrations", "migrator folder path (default ./migrator))")
	migratorCmd.PersistentFlags().StringVarP(&database, "database", "d", "sqlite", "database type ( default sqlite )")

	migratorCmd.AddCommand(migCommands.GenerateCmd)
	migratorCmd.AddCommand(migCommands.UpCmd)
	migratorCmd.AddCommand(migCommands.DownCmd)
	migratorCmd.AddCommand(migCommands.HistoryCmd)

	rootCmd.AddCommand(migratorCmd)
}
