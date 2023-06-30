/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migCommands

import (
	"log"

	"github.com/adharshmk96/stk/pkg/db"
	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/adharshmk96/stk/pkg/migrator/database"
	"github.com/spf13/cobra"
)

// UpCmd represents the mkconfig command
var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "migrate next migrations to database",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootDirectory := cmd.Flag("path").Value.String()
		dbChoice := cmd.Flag("database").Value.String()
		dryRun := cmd.Flag("dry-run").Value.String() == "true"
		numToMigrate := getNumberFromArgs(args, 1)

		conn := db.GetSqliteConnection("migration.db")
		dbRepo := database.NewSqliteRepo(conn)

		log.Println("Generating migration files...")

		config := &migrator.MigratorConfig{
			RootDirectory: rootDirectory,
			Database:      dbChoice,
			NumToMigrate:  numToMigrate,
			DryRun:        dryRun,
			DBRepo:        dbRepo,
		}

		_, err := migrator.MigrateUp(config)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Migrated to database successfully.")

	},
}

func init() {
	UpCmd.Flags().Bool("dry-run", false, "dry run, do not generate files")
}
