/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"fmt"
	"log"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/adharshmk96/stk/pkg/sqlMigrator/dbrepo"
	"github.com/spf13/cobra"
)

func displayRolledBack(rolledBack []*sqlmigrator.MigrationFileEntry) {
	fmt.Printf("\nRolled Back Migrations:\n\n")
	for _, entry := range rolledBack {
		fmt.Println(entry.String())
	}
	fmt.Println("")
}

// DownCmd represents the mkconfig command
var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "Perform backward migration from the files in the migrations folder",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dryRun := cmd.Flag("dry").Value.String() == "true"
		num := getNumberFromArgs(args, 0)

		workDir, dbType, logFile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(workDir, dbType, logFile, dryRun)
		ctx.LoadMigrationEntries()

		dbRepo := dbrepo.SelectDBRepo(dbType)
		migrator := sqlmigrator.NewMigrator(dbRepo)

		rolledBackMigrations, err := migrator.MigrateDown(ctx, num)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = ctx.WriteMigrationEntries()
		if err != nil {
			log.Println("Error writing migration entries:", err)
			return
		}

		displayRolledBack(rolledBackMigrations)
		log.Println("Migrated to database successfully.")

	},
}

func init() {
	DownCmd.Flags().Bool("dry", false, "dry run, do not generate files")
}
