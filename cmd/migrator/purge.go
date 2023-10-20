/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"log"
	"os"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/adharshmk96/stk/pkg/sqlMigrator/dbrepo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PurgeCmd represents the mkconfig command
var PurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "remove all migration files and the migration table from the database",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("purging migrations...")
		rootDirectory := viper.GetString("migrator.workdir")
		err := os.RemoveAll(rootDirectory)
		if err != nil {
			log.Fatal(err)
			return
		}

		dbType := sqlmigrator.SelectDatabase(viper.GetString("migrator.dbtype"))
		dbrepo := dbrepo.SelectDBRepo(dbType)
		err = dbrepo.DeleteMigrationTable()
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("purged migrations successfully.")
	},
}
