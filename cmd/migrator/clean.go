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

func displayCleanedFiles(files []string) {
	fmt.Println("\nCleaned Files:")
	for _, file := range files {
		fmt.Println(file)
	}
	fmt.Println("")
}

// CleanCmd represents the mkconfig command
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all unapplied migration files.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		dryRun := cmd.Flag("dry").Value.String() == "true"

		workDir, dbType, logFile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(workDir, dbType, logFile, dryRun)

		log.Println("Cleaning unapplied migrations...")

		generator := &sqlmigrator.Generator{
			DryRun: dryRun,
		}

		removedFiles, err := generator.Clean(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}

		displayCleanedFiles(removedFiles)
		log.Println("Cleaned migrations successfully.")

	},
}

func init() {
	CleanCmd.Flags().Bool("dry-run", false, "dry run, do not generate files")
}
