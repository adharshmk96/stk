/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"log"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/adharshmk96/stk/pkg/sqlMigrator/dbrepo"
	"github.com/spf13/cobra"
)

// DownCmd represents the mkconfig command
var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "Perform backward migration from the files in the migrations folder",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dryRun := cmd.Flag("dry").Value.String() == "true"
		numToGenerate := getNumberFromArgs(args, 1)

		workDir, dbType, logFile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(workDir, dbType, logFile, dryRun)
		ctx.LoadMigrationEntries()

		dbRepo := dbrepo.SelectDBRepo(dbType)
		migrator := sqlmigrator.NewMigrator(dbRepo)

		_, err := migrator.MigrateDown(ctx, numToGenerate)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Migrated to database successfully.")

	},
}

func init() {
	DownCmd.Flags().Bool("dry", false, "dry run, do not generate files")
}
