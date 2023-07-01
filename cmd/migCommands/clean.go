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
	"github.com/spf13/viper"
)

// CleanCmd represents the mkconfig command
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all unapplied migration files.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootDirectory := viper.GetString("migrator.workdir")
		dbChoice := viper.GetString("migrator.database")

		dryRun := cmd.Flag("dry-run").Value.String() == "true"

		// Select based on the database
		dbType := migrator.SelectDatabase(dbChoice)
		log.Println("selected database: ", dbType)

		extention := migrator.SelectExtention(dbType)
		subDirectory := migrator.SelectSubDirectory(dbType)
		fsRepo := fsrepo.NewFSRepo(filepath.Join(rootDirectory, subDirectory), extention)

		dbRepo := selectDbRepo(dbType)

		log.Println("Cleaning unapplied migrations...")

		config := &migrator.MigratorConfig{
			DryRun: dryRun,

			FSRepo: fsRepo,
			DBRepo: dbRepo,
		}

		_, err := migrator.Clean(config)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Cleaned migrations successfully.")

	},
}

func init() {
	CleanCmd.Flags().Bool("dry-run", false, "dry run, do not generate files")
}
