package sqlmigrator

import (
	"fmt"
	"slices"
)

type MigrationConfig struct {
	NumToApply int
	DBRepo     MigrateDatabase
}

func MigrateUp(ctx *Context, config *MigrationConfig) ([]*MigrationEntry, error) {
	appliedMigrations := []*MigrationEntry{}
	migrationToApply, err := LoadUncommitedMigrations(ctx)
	if err != nil {
		return appliedMigrations, err
	}

	if len(migrationToApply) == 0 {
		fmt.Println("No migrations to apply")
		return appliedMigrations, nil
	}

	num := min(config.NumToApply, len(migrationToApply))
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
		err := config.DBRepo.Exec(upFileContent)
		if err != nil {
			return appliedMigrations, err
		}

		migration.Committed = true
		appliedMigrations = append(appliedMigrations, migration)
	}

	return appliedMigrations, nil
}

func MigrateDown(ctx *Context, config *MigrationConfig) ([]*MigrationEntry, error) {
	rolledBackMigrations := []*MigrationEntry{}
	migrationToApply, err := LoadCommittedMigrations(ctx)
	if err != nil {
		return rolledBackMigrations, err
	}

	if len(migrationToApply) == 0 {
		fmt.Println("No migrations to rollback")
		return rolledBackMigrations, nil
	}

	slices.Reverse(migrationToApply)

	num := min(config.NumToApply, len(migrationToApply))
	if num > 0 {
		migrationToApply = migrationToApply[:num]
	}

	for _, migration := range migrationToApply {
		if ctx.DryRun {
			displayMigration(migration)
			continue
		}

		_, downFileContent := migration.LoadFileContent()

		err := config.DBRepo.Exec(downFileContent)
		if err != nil {
			return rolledBackMigrations, err
		}

		migration.Committed = false
		rolledBackMigrations = append(rolledBackMigrations, migration)
	}

	return rolledBackMigrations, nil
}

func displayMigration(migration *MigrationEntry) {
	fileName := migration.EntryString()
	fmt.Println("up\t:", fileName)
}

func dummyExec(query string) error {
	return nil
}
