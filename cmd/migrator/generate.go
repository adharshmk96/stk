/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"log"
	"strconv"

	"github.com/adharshmk96/stk/pkg/migrator"
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootDirectory := cmd.Flag("path").Value.String()
		database := cmd.Flag("database").Value.String()
		dryRun := cmd.Flag("dry-run").Value.String() == "true"
		numToGenerate := getNumberFromArgs(args, 1)

		log.Println("Generating migration files...")

		config := migrator.GeneratorConfig{
			RootDirectory: rootDirectory,
			Database:      database,
			Name:          migrationName,
			NumToGenerate: numToGenerate,
			DryRun:        dryRun,
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
}
