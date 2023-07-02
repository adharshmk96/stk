/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package cmd

import (
	"github.com/adharshmk96/stk/cmd/migratorCmds"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootFolder string
var database string

// migratorCmd represents the generate command
var migratorCmd = &cobra.Command{
	Use:   "migrator",
	Short: "A database migrator, to generate, run and clean database migrations",
}

func init() {
	migratorCmd.PersistentFlags().StringVarP(&rootFolder, "workdir", "p", "./stk-migrations", "migrator folder (default ./stk-migrations))")
	migratorCmd.PersistentFlags().StringVarP(&database, "database", "d", "sqlite", "database type ( default sqlite )")

	viper.BindPFlag("migrator.workdir", migratorCmd.PersistentFlags().Lookup("workdir"))
	viper.BindPFlag("migrator.database", migratorCmd.PersistentFlags().Lookup("database"))

	migratorCmd.AddCommand(migratorCmds.GenerateCmd)
	migratorCmd.AddCommand(migratorCmds.UpCmd)
	migratorCmd.AddCommand(migratorCmds.DownCmd)
	migratorCmd.AddCommand(migratorCmds.CleanCmd)
	migratorCmd.AddCommand(migratorCmds.HistoryCmd)

	rootCmd.AddCommand(migratorCmd)

}
