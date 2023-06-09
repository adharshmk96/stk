/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migratorCmds

import (
	"log"

	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// historyCmd represents the mkconfig command
var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View the migration history of the database.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		dbChoice := viper.GetString("migrator.database")

		// Select based on the database
		dbType := migrator.SelectDatabase(dbChoice)

		dbRepo := selectDbRepo(dbType)

		err := migrator.MigrationHistory(dbRepo)
		if err != nil {
			log.Fatal(err)
			return
		}

	},
}
