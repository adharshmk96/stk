package sqlmigrator

import "fmt"

func MigrateUp(ctx *Context, numToApply int) ([]*MigrationEntry, error) {
	appliedMigrations := []*MigrationEntry{}
	migrationToApply, err := LoadUncommitedMigrations(ctx)
	if err != nil {
		return appliedMigrations, err
	}

	if len(migrationToApply) == 0 {
		fmt.Println("No migrations to apply")
		return appliedMigrations, nil
	}

	num := min(numToApply, len(migrationToApply))
	if num > 0 {
		migrationToApply = migrationToApply[:num]
	}

	for _, migration := range migrationToApply {
		if ctx.DryRun {
			displayMigration(migration)
			continue
		}

		upFileContent, _ := migration.LoadFileContent()

		// TODO: replace with db stuff
		err := dummyExec(upFileContent)
		if err != nil {
			return appliedMigrations, err
		}

		migration.Committed = true
		appliedMigrations = append(appliedMigrations, migration)
	}

	return appliedMigrations, nil
}

func displayMigration(migration *MigrationEntry) {
	fileName := migration.EntryString()
	fmt.Println("up\t:", fileName)
}

func dummyExec(query string) error {
	return nil
}
