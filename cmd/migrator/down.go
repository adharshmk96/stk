/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"log"
	"path/filepath"

	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/adharshmk96/stk/pkg/migrator/fsrepo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DownCmd represents the mkconfig command
var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "Perform backward migration from the files in the migrations folder",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootDirectory := viper.GetString("migrator.workdir")
		dbChoice := viper.GetString("migrator.database")

		dryRun := cmd.Flag("dry-run").Value.String() == "true"

		numToMigrate := getNumberFromArgs(args, 0)

		// Select based on the database
		dbType := migrator.SelectDatabase(dbChoice)
		log.Println("selected database: ", dbType)

		extention := migrator.SelectExtention(dbType)
		subDirectory := migrator.SelectSubDirectory(dbType)
		fsRepo := fsrepo.NewFSRepo(filepath.Join(rootDirectory, subDirectory), extention)

		dbRepo := selectDbRepo(dbType)

		log.Println("Applying migrations down...")

		config := &migrator.MigratorConfig{
			NumToMigrate: numToMigrate,
			DryRun:       dryRun,

			FSRepo: fsRepo,
			DBRepo: dbRepo,
		}

		_, err := migrator.MigrateDown(config)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Migrated to database successfully.")

	},
}

func init() {
	DownCmd.Flags().Bool("dry-run", false, "dry run, do not generate files")
}
