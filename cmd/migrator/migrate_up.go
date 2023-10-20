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

// UpCmd represents the mkconfig command
var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "Perform forward migration from the files in the migrations folder",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		dryRun := cmd.Flag("dry").Value.String() == "true"
		num := getNumberFromArgs(args, 1)

		workDir, dbType, logFile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(workDir, dbType, logFile, dryRun)
		ctx.LoadMigrationEntries()

		dbRepo := dbrepo.SelectDBRepo(dbType)
		migrator := sqlmigrator.NewMigrator(dbRepo)

		_, err := migrator.MigrateUp(ctx, num)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Migrated to database successfully.")

	},
}

func init() {
	UpCmd.Flags().Bool("dry", false, "dry run, do not generate files")
}
