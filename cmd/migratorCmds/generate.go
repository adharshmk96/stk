/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migratorCmds

import (
	"fmt"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/spf13/cobra"
)

var migrationName string

func displayContext(ctx *sqlmigrator.Context) {
	fmt.Println("Work Directory: ", ctx.WorkDir)
	fmt.Println("Log File: ", ctx.LogFile)
	fmt.Println("Database: ", ctx.Database)
}

func displayGenerator(generator *sqlmigrator.Generator) {
	fmt.Println("Migration Name: ", generator.Name)
	fmt.Println("Number of Migrations: ", generator.NumToGenerate)
	fmt.Println("Dry Run: ", generator.DryRun)
	fmt.Println("Fill: ", generator.Fill)
}

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate migration files.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		dryRun := cmd.Flag("dry").Value.String() == "true"
		fill := cmd.Flag("fill").Value.String() == "true"

		numToGenerate := getNumberFromArgs(args, 1)

		ctx := sqlmigrator.NewMigratorContext(dryRun)
		displayContext(ctx)
		generator := sqlmigrator.NewGenerator(migrationName, numToGenerate, dryRun, fill)
		displayGenerator(generator)

		err := generator.Generate(ctx)
		if err != nil {
			return
		}

		fmt.Println("Generated migrations successfully.")

	},
}

func init() {
	GenerateCmd.Flags().StringVarP(&migrationName, "name", "n", "", "migration name")
	GenerateCmd.Flags().Bool("dry", false, "dry run, do not generate files")
	GenerateCmd.Flags().Bool("fill", false, "fill the created files with sample content")

}
