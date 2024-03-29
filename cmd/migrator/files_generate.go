/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"fmt"
	"log"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/spf13/cobra"
)

var migrationName string

func displayContextAndConfig(ctx *sqlmigrator.Context, generator *sqlmigrator.Generator) {
	labels := []string{"Work Directory", "Log File", "Database", "Name", "Files", "Dry Run", "Fill"}

	maxLen := 0
	for _, label := range labels {
		if len(label) > maxLen {
			maxLen = len(label)
		}
	}

	fmt.Printf("%-*s :%v\n", maxLen, labels[0], ctx.WorkDir)
	fmt.Printf("%-*s :%v\n", maxLen, labels[1], ctx.LogFile)
	fmt.Printf("%-*s :%v\n", maxLen, labels[2], ctx.Database)
	fmt.Printf("%-*s :%v\n", maxLen, labels[3], generator.Name)
	fmt.Printf("%-*s :%v\n", maxLen, labels[4], generator.NumToGenerate)
	fmt.Printf("%-*s :%v\n", maxLen, labels[5], generator.DryRun)
	fmt.Printf("%-*s :%v\n", maxLen, labels[6], generator.Fill)

}

func displayGeneratedFiles(files []string) {
	fmt.Println("\ngenerated files:")
	for _, file := range files {
		fmt.Println(file)
	}
	fmt.Println("")
}

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate migration files.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		dryRun := cmd.Flag("dry").Value.String() == "true"
		fill := cmd.Flag("fill").Value.String() == "true"

		numToGenerate := getNumberFromArgs(args, 1)

		workDir, dbType, logFile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(workDir, dbType, logFile, dryRun)
		ctx.LoadMigrationEntries()
		generator := sqlmigrator.NewGenerator(migrationName, numToGenerate, fill)
		displayContextAndConfig(ctx, generator)

		log.Println("generating migrations...")
		generatedFiles, err := generator.Generate(ctx)
		if err != nil {
			log.Println("error generating migrations:", err)
			return
		}
		displayGeneratedFiles(generatedFiles)

		err = ctx.WriteMigrationEntries()
		if err != nil {
			log.Println("error writing migration entries:", err)
			return
		}
		log.Println("generated migrations successfully.")

	},
}

func init() {
	GenerateCmd.Flags().StringVarP(&migrationName, "name", "n", "", "migration name")
	GenerateCmd.Flags().Bool("dry", false, "dry run, do not generate files")
	GenerateCmd.Flags().Bool("fill", false, "fill the created files with sample content")

}
