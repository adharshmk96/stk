/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"fmt"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/spf13/cobra"
)

var migrationName string

func displayContext(ctx *sqlmigrator.Context) {
	fmt.Println("Work Directory\t:", ctx.WorkDir)
	fmt.Println("Log File\t:", ctx.LogFile)
	fmt.Println("Database:\t:", ctx.Database)
}

func displayGenerator(generator *sqlmigrator.Generator) {
	fmt.Println("Name\t:", generator.Name)
	fmt.Println("Files\t:", generator.NumToGenerate)
	fmt.Println("Dry Run\t:", generator.DryRun)
	fmt.Println("Fill\t:", generator.Fill)
}

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate migration files.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		dryRun := cmd.Flag("dry").Value.String() == "true"
		fill := cmd.Flag("fill").Value.String() == "true"

		numToGenerate := getNumberFromArgs(args, 1)

		workDir, dbType, logFile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(workDir, dbType, logFile, dryRun)
		displayContext(ctx)
		generator := sqlmigrator.NewGenerator(migrationName, numToGenerate, dryRun, fill)
		displayGenerator(generator)

		fmt.Println("Generating migrations...")
		err := generator.Generate(ctx)
		if err != nil {
			fmt.Println("Error generating migrations:", err)
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
