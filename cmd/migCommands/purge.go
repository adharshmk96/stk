/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migCommands

import (
	"log"
	"path/filepath"

	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/adharshmk96/stk/pkg/migrator/fsrepo"
	"github.com/spf13/cobra"
)

// PurgeCmd represents the mkconfig command
var PurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Remove all migration files and the migration table from the database",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootDirectory := cmd.Flag("path").Value.String()
		dbChoice := cmd.Flag("database").Value.String()

		// Select based on the database
		dbType := migrator.SelectDatabase(dbChoice)
		log.Println("selected database: ", dbType)

		extention := migrator.SelectExtention(dbType)
		subDirectory := migrator.SelectSubDirectory(dbType)
		fsRepo := fsrepo.NewFSRepo(filepath.Join(rootDirectory, subDirectory), extention)

		dbRepo := selectDbRepo(dbType)

		log.Println("Removing migration table...")

		err := dbRepo.DeleteMigrationTable()
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Removed migration table successfully.")

		log.Println("Removing migration files")

		err = fsRepo.DeleteMigrationDirectory()
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Removed migration files successfully.")

	},
}

func init() {
	PurgeCmd.Flags().Bool("dry-run", false, "dry run, do not generate files")
}
