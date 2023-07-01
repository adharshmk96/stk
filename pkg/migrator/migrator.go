package migrator

import (
	"log"
	"sort"
)

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

	if lastAppliedMigration != nil {
		log.Println("last applied migration : ", lastAppliedMigration)
	} else {
		log.Println("last applied migration : ", "Empty")
	}

	// Find the next migrations to apply
	migrationsToApply := CalculateUpMigrationsToApply(lastAppliedMigration, migrations, config.NumToMigrate)

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

	reverseMigrationList(migrations)
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
		// Find the index of the element just greater than num
		startIdx := sort.Search(len(migrations), func(i int) bool { return migrations[i].Number > lastMigration.Number })

		// If the index + n is out of bounds
		endInx := min(startIdx+numberToMigrate, len(migrations))

		result := migrations[startIdx:endInx]
		return result
	}

	if lastMigration.Type == MigrationDown {
		// Find the index of the element just greater than num
		startIdx := sort.Search(len(migrations), func(i int) bool { return migrations[i].Number > lastMigration.Number })

		startIdx = max(startIdx-1, 0)
		endInx := min(startIdx+numberToMigrate, len(migrations))

		result := migrations[startIdx:endInx]
		return result
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
		// Find the next migration value than last applied migration
		index := sort.Search(len(migrations), func(i int) bool { return migrations[i].Number <= lastMigration.Number })

		// If the index - n is less than 0, return an error
		endIdx := min(index+numberToMigrate, len(migrations))

		// Otherwise, slice the list to get the last n numbers smaller than num
		result := migrations[index:endIdx]

		return result
	}

	if lastMigration.Type == MigrationDown {
		// Find the index of the element just greater or equal to num
		index := sort.Search(len(migrations), func(i int) bool { return migrations[i].Number <= lastMigration.Number }) + 1

		// If the index - n is less than 0, return an error
		endIdx := min(index+numberToMigrate, len(migrations))

		// Otherwise, slice the list to get the last n numbers smaller than num
		result := migrations[index:endIdx]

		return result
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
