package migrator

import "log"

type MigratorConfig struct {
	NumToMigrate int
	DryRun       bool
	DBRepo       DatabaseRepo
	FSRepo       FileRepo
}

func MigrateUp(config *MigratorConfig) ([]*Migration, error) {

	fsRepo := config.FSRepo
	dbRepo := config.DBRepo

	if config.DBRepo == nil {
		log.Fatalf("database is not initialized")
		return nil, ErrDatabaseNotInitialized
	}

	// Read last applied migration from database
	lastAppliedMigration, err := dbRepo.LoadLastAppliedMigration()
	if err != nil {
		return nil, err
	}

	// Read and parse all migrations from directory
	migrations, err := fsRepo.LoadMigrationsFromFile(MigrationUp)
	if err != nil {
		return nil, err
	}

	log.Println("loaded migrations: ")
	for _, v := range migrations {
		log.Println(" - ", v.Number, v.Name)
	}

	log.Println("last applied migration : ", lastAppliedMigration)

	// Find the next migrations to apply
	migrationsToApply := findUpMigrationsToApply(lastAppliedMigration, migrations, config.NumToMigrate)

	// Read migration queries from files
	for _, migration := range migrationsToApply {
		err := fsRepo.LoadMigrationQuery(migration)
		if err != nil {
			return nil, err
		}

		// Apply migrations and add entries to database
		err = applyMigration(config, migration)
		if err != nil {
			return nil, err
		}
	}

	return migrationsToApply, nil

}

// TODO: apply all as one transaction ?. if one fails, rollback all
func applyMigration(config *MigratorConfig, migration *Migration) error {
	if config.DryRun {
		log.Println("dry run: ", migration.Number, migration.Name)
		return nil
	}

	log.Println("applying migration: ", migration.Number, migration.Name)
	err := config.DBRepo.ApplyMigration(migration)
	if err != nil {
		log.Fatalln("error applying migration: ", migration.Number, err)
		return err
	}

	return nil
}

func findUpMigrationsToApply(lastMigration *Migration, migrations []*Migration, numberToMigrate int) []*Migration {
	var startIdx, lastIdx int
	if lastMigration == nil {
		startIdx = 0
		lastIdx = min(startIdx+numberToMigrate, len(migrations))
	} else {
		idx := findNextMigrationIndex(migrations, lastMigration.Number)
		if lastMigration.Type == MigrationDown {
			startIdx = max(idx-1, 0)
			lastIdx = min(startIdx+numberToMigrate, len(migrations))
		} else {
			startIdx = max(idx, 0)
			lastIdx = min(startIdx+numberToMigrate, len(migrations))
		}
	}
	return migrations[startIdx:lastIdx]
}

func findNextMigrationIndex(migrations []*Migration, number int) int {
	for i, v := range migrations {
		if v.Number > number {
			return i
		}
	}
	return -1
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
