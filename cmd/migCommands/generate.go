/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migCommands

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/adharshmk96/stk/pkg/migrator/fsrepo"
	"github.com/spf13/cobra"
)

var migrationName string

func getNumberFromArgs(args []string, defaultValue int) int {
	if len(args) == 0 {
		return defaultValue
	}
	num, err := strconv.Atoi(args[0])
	if err != nil {
		return defaultValue
	}
	return num
}

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate migration files",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootDirectory := cmd.Flag("path").Value.String()
		dbChoice := cmd.Flag("database").Value.String()

		dryRun := cmd.Flag("dry-run").Value.String() == "true"
		fill := cmd.Flag("fill").Value.String() == "true"

		numToGenerate := getNumberFromArgs(args, 1)

		// Select based on the database
		database := migrator.SelectDatabase(dbChoice)
		log.Println("selected database: ", database)

		extention := migrator.GetExtention(database)
		subDirectory := migrator.SelectSubDirectory(database)
		fsRepo := fsrepo.NewFSRepo(filepath.Join(rootDirectory, subDirectory), extention)

		log.Println("Generating migration files...")

		config := migrator.GeneratorConfig{
			Database:      database,
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
