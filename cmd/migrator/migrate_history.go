/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package migrator

import (
	"fmt"
	"log"
	"strings"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/adharshmk96/stk/pkg/sqlMigrator/dbrepo"
	"github.com/spf13/cobra"
)

func displayMigrationHistory(history []*sqlmigrator.MigrationDBEntry) {
	// Define max column widths
	numberWidth := len("number")
	nameWidth := len("name")
	directionWidth := len("direction")
	createdWidth := len("created")

	// Find the maximum width needed for each column based on data
	for _, entry := range history {
		if len(entry.Name) > nameWidth {
			nameWidth = len(entry.Name)
		}
		if len(entry.Direction) > directionWidth {
			directionWidth = len(entry.Direction)
		}
		if len(entry.Created.String()) > createdWidth {
			createdWidth = len(entry.Created.String())
		}
	}

	// Print header
	fmt.Println("Migration History")
	fmt.Println("-----------------")
	fmt.Printf("| %-"+fmt.Sprintf("%d", numberWidth)+"s | %-"+fmt.Sprintf("%d", nameWidth)+"s | %-"+fmt.Sprintf("%d", directionWidth)+"s | %-"+fmt.Sprintf("%d", createdWidth)+"s |\n", "number", "name", "direction", "created")
	fmt.Printf("| %-"+fmt.Sprintf("%d", numberWidth)+"s | %-"+fmt.Sprintf("%d", nameWidth)+"s | %-"+fmt.Sprintf("%d", directionWidth)+"s | %-"+fmt.Sprintf("%d", createdWidth)+"s |\n", strings.Repeat("-", numberWidth), strings.Repeat("-", nameWidth), strings.Repeat("-", directionWidth), strings.Repeat("-", createdWidth))

	// Print entries
	for _, entry := range history {
		fmt.Printf("| %-"+fmt.Sprintf("%d", numberWidth)+"d | %-"+fmt.Sprintf("%d", nameWidth)+"s | %-"+fmt.Sprintf("%d", directionWidth)+"s | %-"+fmt.Sprintf("%d", createdWidth)+"s |\n", entry.Number, entry.Name, entry.Direction, entry.Created.String())
	}
}

// historyCmd represents the mkconfig command
var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View the migration history of the database.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		workDir, dbType, logFile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(workDir, dbType, logFile, false)

		dbRepo := dbrepo.SelectDBRepo(dbType)
		migrator := sqlmigrator.NewMigrator(dbRepo)

		history, err := migrator.MigrationHistory(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}

		displayMigrationHistory(history)

	},
}
