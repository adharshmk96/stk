/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PurgeCmd represents the mkconfig command
var PurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Remove all migration files and the migration table from the database",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootDirectory := viper.GetString("migrator.workdir")
		err := os.RemoveAll(rootDirectory)
		if err != nil {
			log.Fatal(err)
			return
		}

		// TODO: Remove the migration table from the database

		log.Println("Purged migrations successfully.")
	},
}

func init() {
	PurgeCmd.Flags().Bool("dry-run", false, "dry run, do not generate files")
}
