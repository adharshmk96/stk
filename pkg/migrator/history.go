package migrator

import (
	"fmt"
	"strings"
)

func MigrationHistory(dbRepo DatabaseRepo) error {
	migrations, err := dbRepo.LoadMigrations()
	if err != nil {
		return err
	}

	printTable(migrations)

	return nil
}

func printTable(data []*Migration) {
	header := []string{"Number", "Name", "Type", "Created"}
	entries := make([][]string, len(data)+1)
	columns := len(header)

	entries[0] = header

	for i, entry := range data {
		entries[i+1] = []string{fmt.Sprintf("%06d", entry.Number), entry.Name, string(entry.Type), entry.Created.Format("2006-01-02")}
	}

	maxWidths := make([]int, columns)

	for _, row := range entries {
		for i, item := range row {
			if len(item) > maxWidths[i] {
				maxWidths[i] = len(item)
			}
		}
	}

	for _, row := range entries {
		sb := strings.Builder{}

		for i, item := range row {
			sb.WriteString(item + strings.Repeat(" ", maxWidths[i]-len(item)+2))
		}

		fmt.Println(sb.String())
	}
}
