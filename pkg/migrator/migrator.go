package migrator

import (
	"log"
)

type MigratorConfig struct {
	NumToMigrate int
	DryRun       bool
	DBRepo       DatabaseRepo
	FSRepo       FileRepo
}

func MigrateUp(config *MigratorConfig) ([]*Migration, error) {
	migrations, lastAppliedMigration, err := LoadMigrations(config, MigrationUp)
	if err != nil {
		return nil, err
	}

	// Find the next migrations to apply
	migrationsToApply := CalculateUpMigrationsToApply(lastAppliedMigration, migrations, config.NumToMigrate)
	if len(migrationsToApply) == 0 {
		log.Println("No migrations to apply")
		return nil, nil
	}

	// Apply migrations
	err = ApplyMigrations(config, migrationsToApply)
	if err != nil {
		return nil, err
	}

	return migrationsToApply, nil
}

func MigrateDown(config *MigratorConfig) ([]*Migration, error) {
	migrations, lastAppliedMigration, err := LoadMigrations(config, MigrationDown)
	if err != nil {
		return nil, err
	}

	// Sort migrations in reverse order
	reverseMigrationList(migrations)
	// Find the next migrations to apply
	migrationsToApply := CalculateDownMigrationsToApply(lastAppliedMigration, migrations, config.NumToMigrate)
	if len(migrationsToApply) == 0 {
		log.Println("No migrations to apply")
		return nil, nil
	}

	// Apply migrations
	err = ApplyMigrations(config, migrationsToApply)
	if err != nil {
		return nil, err
	}

	return migrationsToApply, nil
}

func LoadMigrations(config *MigratorConfig, migrationType MigrationType) ([]*Migration, *Migration, error) {
	fsRepo := config.FSRepo
	dbRepo := config.DBRepo

	if config.DBRepo == nil {
		log.Fatalf("database is not initialized")
		return nil, nil, ErrDatabaseNotInitialized
	}

	// Read last applied migration from database
	lastAppliedMigration, err := dbRepo.LoadLastAppliedMigration()
	if err != nil {
		return nil, nil, err
	}

	// Read and parse all migrations from directory
	migrations, err := fsRepo.LoadMigrationsFromFile(migrationType)
	if err != nil {
		return nil, nil, err
	}

	return migrations, lastAppliedMigration, nil
}

func ApplyMigrations(config *MigratorConfig, migrationsToApply []*Migration) error {
	fsRepo := config.FSRepo

	// Read migration queries from files
	for _, migration := range migrationsToApply {
		err := fsRepo.LoadMigrationQuery(migration)
		if err != nil {
			return err
		}

		// Apply migrations and add entries to database
		err = ApplyMigration(config, migration)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: apply all as one transaction ?. if one fails, rollback all
func ApplyMigration(config *MigratorConfig, migration *Migration) error {
	if config.DryRun {
		log.Println("dry run: ", migration.Number, migration.Name)
		return nil
	}

	log.Println("applying migration: ", migration.Number, migration.Name)
	err := config.DBRepo.ApplyMigration(migration)
	if err != nil {
		log.Println("error applying migration: ", migration.Number, err)
		return err
	}

	return nil
}

func CalculateUpMigrationsToApply(lastMigration *Migration, migrations []*Migration, numberToMigrate int) []*Migration {
	if numberToMigrate == 0 {
		numberToMigrate = len(migrations)
	}

	if len(migrations) == 0 {
		return []*Migration{}
	}

	if lastMigration == nil {
		endIdx := min(numberToMigrate, len(migrations))
		return migrations[:endIdx]
	}

	if lastMigration.Type == MigrationUp {
		for i, v := range migrations {
			if v.Number > lastMigration.Number {
				endIdx := min(i+numberToMigrate, len(migrations))
				return migrations[i:endIdx]
			}
		}
	}

	if lastMigration.Type == MigrationDown {
		for i, v := range migrations {
			if v.Number >= lastMigration.Number {
				endIdx := min(i+numberToMigrate, len(migrations))
				return migrations[i:endIdx]
			}
		}
	}

	return []*Migration{}

}

// It will accept a reversed list of migrations
func CalculateDownMigrationsToApply(lastMigration *Migration, migrations []*Migration, numberToMigrate int) []*Migration {
	if numberToMigrate == 0 {
		numberToMigrate = len(migrations)
	}

	if len(migrations) == 0 {
		return []*Migration{}
	}

	if lastMigration == nil {
		endIdx := min(numberToMigrate, len(migrations))
		return migrations[:endIdx]
	}

	if lastMigration.Type == MigrationUp {
		for i, v := range migrations {
			if v.Number <= lastMigration.Number {
				endIdx := min(i+numberToMigrate, len(migrations))
				return migrations[i:endIdx]
			}
		}
	}

	if lastMigration.Type == MigrationDown {
		for i, v := range migrations {
			if v.Number < lastMigration.Number {
				endIdx := min(i+numberToMigrate, len(migrations))
				return migrations[i:endIdx]
			}
		}
	}

	return []*Migration{}
}

func reverseMigrationList(migrations []*Migration) []*Migration {
	for i := len(migrations)/2 - 1; i >= 0; i-- {
		opp := len(migrations) - 1 - i
		migrations[i], migrations[opp] = migrations[opp], migrations[i]
	}
	return migrations
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
