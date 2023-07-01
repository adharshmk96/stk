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

var migrationName string

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate migration files.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootDirectory := viper.GetString("migrator.workdir")
		dbChoice := viper.GetString("migrator.database")

		dryRun := cmd.Flag("dry-run").Value.String() == "true"
		fill := cmd.Flag("fill").Value.String() == "true"

		numToGenerate := getNumberFromArgs(args, 1)

		// Select based on the dbType
		dbType := migrator.SelectDatabase(dbChoice)
		log.Println("selected database: ", dbType)

		extention := migrator.SelectExtention(dbType)
		subDirectory := migrator.SelectSubDirectory(dbType)
		fsRepo := fsrepo.NewFSRepo(filepath.Join(rootDirectory, subDirectory), extention)

		log.Println("Generating migration files...")

		config := migrator.GeneratorConfig{
			Name:          migrationName,
			NumToGenerate: numToGenerate,
			DryRun:        dryRun,
			Fill:          fill,
			FSRepo:        fsRepo,
		}

		err := migrator.Generate(config)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Migration files generated successfully.")
	},
}

func init() {
	GenerateCmd.Flags().StringVarP(&migrationName, "name", "n", "", "migration name")
	GenerateCmd.Flags().Bool("dry-run", false, "dry run, do not generate files")
	GenerateCmd.Flags().Bool("fill", false, "fill the created files with sample content")

}
