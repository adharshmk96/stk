package sqlmigrator

import (
	"fmt"
	"slices"
)

type migrator struct {
	DBRepo DBRepo
}

func NewMigrator(dbRepo DBRepo) *migrator {
	return &migrator{
		DBRepo: dbRepo,
	}
}

func (m *migrator) MigrateUp(ctx *Context, num int) ([]*MigrationFileEntry, error) {
	appliedMigrations := []*MigrationFileEntry{}
	migrationToApply, err := LoadUncommitedMigrations(ctx)
	if err != nil {
		return appliedMigrations, err
	}

	if len(migrationToApply) == 0 {
		fmt.Println("No migrations to apply")
		return appliedMigrations, nil
	}

	num = min(num, len(migrationToApply))
	if num > 0 {
		migrationToApply = migrationToApply[:num]
	}

	for _, migration := range migrationToApply {
		if ctx.DryRun {
			displayMigration(migration)
			continue
		}

		upFileContent, _ := migration.LoadFileContent()

		err := m.DBRepo.Exec(upFileContent)
		if err != nil {
			return appliedMigrations, err
		}

		migration.Committed = true
		appliedMigrations = append(appliedMigrations, migration)

		// commit to db migration table
		dbEntry := &MigrationDBEntry{
			Name:      migration.String(),
			Direction: "up",
		}

		err = m.DBRepo.PushHistory(dbEntry)
		if err != nil {
			return appliedMigrations, err
		}
	}

	return appliedMigrations, nil
}

func (m *migrator) MigrateDown(ctx *Context, num int) ([]*MigrationFileEntry, error) {
	rolledBackMigrations := []*MigrationFileEntry{}
	migrationToApply, err := LoadCommittedMigrations(ctx)
	if err != nil {
		return rolledBackMigrations, err
	}

	if len(migrationToApply) == 0 {
		fmt.Println("No migrations to rollback")
		return rolledBackMigrations, nil
	}

	slices.Reverse(migrationToApply)

	num = min(num, len(migrationToApply))
	if num > 0 {
		migrationToApply = migrationToApply[:num]
	}

	for _, migration := range migrationToApply {
		if ctx.DryRun {
			displayMigration(migration)
			continue
		}

		_, downFileContent := migration.LoadFileContent()

		err := m.DBRepo.Exec(downFileContent)
		if err != nil {
			return rolledBackMigrations, err
		}

		migration.Committed = false
		rolledBackMigrations = append(rolledBackMigrations, migration)

		// commit to db migration table
		dbEntry := &MigrationDBEntry{
			Name:      migration.String(),
			Direction: "down",
		}

		err = m.DBRepo.PushHistory(dbEntry)
		if err != nil {
			return rolledBackMigrations, err
		}
	}

	return rolledBackMigrations, nil
}

func (m *migrator) MigrationHistory(ctx *Context) ([]*MigrationDBEntry, error) {
	return m.DBRepo.LoadHistory()
}

func displayMigration(migration *MigrationFileEntry) {
	fileName := migration.EntryString()
	fmt.Println("up\t:", fileName)
}
