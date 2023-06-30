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
	migrationsToApply := FindUpMigrationsToApply(lastAppliedMigration, migrations, config.NumToMigrate)

	if len(migrationsToApply) == 0 {
		log.Println("No migrations to apply")
		return nil, nil
	}

	// Read migration queries from files
	for _, migration := range migrationsToApply {
		err := fsRepo.LoadMigrationQuery(migration)
		if err != nil {
			return nil, err
		}

		// Apply migrations and add entries to database
		err = ApplyMigration(config, migration)
		if err != nil {
			return nil, err
		}
	}

	return migrationsToApply, nil

}

func MigrateDown(config *MigratorConfig) ([]*Migration, error) {

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
	migrations, err := fsRepo.LoadMigrationsFromFile(MigrationDown)
	if err != nil {
		return nil, err
	}

	log.Println("loaded migrations: ")
	for _, v := range migrations {
		log.Println(" - ", v.Number, v.Name)
	}

	log.Println("last applied migration : ", lastAppliedMigration)

	// Find the next migrations to apply
	migrationsToApply := CalculateDownMigrationsToApply(lastAppliedMigration, migrations, config.NumToMigrate)

	if len(migrationsToApply) == 0 {
		log.Println("No migrations to apply")
		return nil, nil
	}

	// Read migration queries from files
	for _, migration := range migrationsToApply {
		err := fsRepo.LoadMigrationQuery(migration)
		if err != nil {
			return nil, err
		}

		// Apply migrations and add entries to database
		err = ApplyMigration(config, migration)
		if err != nil {
			return nil, err
		}
	}

	return migrationsToApply, nil

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

func getStartIdx(lastMigration *Migration, migrations []*Migration) int {
	if lastMigration == nil {
		return 0
	}

	idx := findNextMigrationNumberIndex(migrations, lastMigration.Number)
	if idx < 0 {
		return -1
	}
	if lastMigration.Type == MigrationDown {
		return max(idx-1, 0)
	}
	return max(idx, 0)
}

func getLastIdx(startIdx int, numberToMigrate int, totalMigrations int) int {
	if numberToMigrate == 0 || startIdx == -1 {
		return totalMigrations
	}
	return min(startIdx+numberToMigrate, totalMigrations)
}

func FindUpMigrationsToApply(lastMigration *Migration, migrations []*Migration, numberToMigrate int) []*Migration {
	startIdx := getStartIdx(lastMigration, migrations)
	lastIdx := getLastIdx(startIdx, numberToMigrate, len(migrations))
	if startIdx == -1 {
		return []*Migration{}
	}
	return migrations[startIdx:lastIdx]
}

func CalculateUpMigrationsToApply(lastMigration *Migration, migrations []*Migration, numberToMigrate int) []*Migration {
	if lastMigration == nil {
		return migrations[:numberToMigrate]
	}

	idx := findNextMigrationNumberIndex(migrations, lastMigration.Number)
	if idx < 0 {
		return []*Migration{}
	}

	var startIdx, endIdx int
	if lastMigration.Type == MigrationDown {
		startIdx = max(idx-1, 0)
		if numberToMigrate <= 0 {
			endIdx = len(migrations)
		} else {
			endIdx = min(startIdx+numberToMigrate, len(migrations))
		}

	} else {
		startIdx = max(idx, 0)
		if numberToMigrate <= 0 {
			endIdx = len(migrations)
		} else {
			endIdx = min(startIdx+numberToMigrate, len(migrations))
		}
	}

	return migrations[startIdx:endIdx]
}

func CalculateDownMigrationsToApply(lastMigration *Migration, migrations []*Migration, numberToMigrate int) []*Migration {
	if lastMigration == nil {
		return []*Migration{}
	}

	idx := findNextMigrationNumberIndex(migrations, lastMigration.Number)
	if idx < 0 {
		return []*Migration{}
	}

	var startIdx, endIdx int
	if lastMigration.Type == MigrationDown {
		endIdx = max(idx-1, 0)
		if numberToMigrate <= 0 {
			startIdx = 0
		} else {
			startIdx = max(idx-numberToMigrate-1, 0)
		}

	} else {
		endIdx = max(idx, 0)
		if numberToMigrate <= 0 {
			startIdx = 0
		} else {
			startIdx = max(endIdx-numberToMigrate, 0)
		}
	}

	downMigrations := migrations[startIdx:endIdx]

	// reverse order
	for i, j := 0, len(downMigrations)-1; i < j; i, j = i+1, j-1 {
		downMigrations[i], downMigrations[j] = downMigrations[j], downMigrations[i]
	}

	return downMigrations

}

func findNextMigrationNumberIndex(migrations []*Migration, number int) int {
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
