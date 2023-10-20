/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"log"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/spf13/cobra"
)

// historyCmd represents the mkconfig command
var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View the migration history of the database.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		workDir, dbType, logFile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(workDir, dbType, logFile, false)
		dbRepo := sqlmigrator.SelectDBRepo(dbType, "path")
		migrator := sqlmigrator.NewMigrator(dbRepo)

		_, err := migrator.MigrationHistory(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Migrated to database successfully.")

	},
}
