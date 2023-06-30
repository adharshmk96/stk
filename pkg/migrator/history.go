package migrator

import (
	"fmt"
	"time"
)

func MigrationHistory(dbRepo DatabaseRepo) error {
	migrations, err := dbRepo.LoadMigrations()
	if err != nil {
		return err
	}

	displayMigrations(migrations)

	return nil
}

func displayMigrations(migrations []*Migration) {

	if len(migrations) == 0 {
		fmt.Println("No migrations found...")
		return
	}

	// Print headers
	fmt.Printf("%-10s %-20s %-10s %-30s\n", "Number", "Name", "Type", "Created")

	// Print each migration
	for _, migration := range migrations {
		fmt.Printf("%-10d %-20s %-10s %-30s\n",
			migration.Number,
			migration.Name,
			string(migration.Type),
			migration.Created.Format(time.RFC3339))
	}
}
