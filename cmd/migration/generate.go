/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package migration

import (
	"log"

	"github.com/adharshmk96/stk/internal/migrator"
	"github.com/spf13/cobra"
)

var migrationName string

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate migration files",
	Run: func(cmd *cobra.Command, args []string) {
		rootDirectory := cmd.Flag("path").Value.String()
		database := cmd.Flag("database").Value.String()

		log.Println("Generating migration files...")

		config := migrator.GeneratorConfig{
			RootDirectory: rootDirectory,
			Database:      database,
			Name:          migrationName,
			NumToGenerate: 1,
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
}
