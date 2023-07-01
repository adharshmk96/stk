package migrator

import (
	"log"
)

func Clean(config *MigratorConfig) ([]*Migration, error) {

	fsRepo := config.FSRepo
	dbRepo := config.DBRepo

	// Read and parse all migrations from directory
	upMigrations, err := fsRepo.LoadMigrationsFromFile(MigrationUp)
	if err != nil {
		return nil, err
	}
	downMigrations, err := fsRepo.LoadMigrationsFromFile(MigrationDown)
	if err != nil {
		return nil, err
	}

	// Read last applied migration from database
	lastAppliedMigration, err := dbRepo.LoadLastAppliedMigration()
	if err != nil {
		return nil, err
	}

	// Find the next migrations to apply
	migrationsToApply := CalculateUpMigrationsToApply(lastAppliedMigration, upMigrations, config.NumToMigrate)
	if len(migrationsToApply) == 0 {
		log.Println("No migrations to clean")
		return nil, nil
	}

	downMigrations = filterMigrations(downMigrations, migrationsToApply)

	if config.DryRun {
		for _, migration := range migrationsToApply {
			log.Println("Cleaning migration", migration.Path)
		}
		for _, migration := range downMigrations {
			log.Println("Cleaning migration", migration.Path)
		}
		return migrationsToApply, nil
	}

	// Clear migrations
	for _, migration := range migrationsToApply {
		log.Println("Cleaning migration", migration.Path)
		err = config.FSRepo.DeleteMigrationFile(migration)
		if err != nil {
			return nil, err
		}
	}
	for _, downMigration := range downMigrations {
		log.Println("Cleaning migration", downMigration.Path)
		err = config.FSRepo.DeleteMigrationFile(downMigration)
		if err != nil {
			return nil, err
		}
	}

	return migrationsToApply, nil
}

func filterMigrations(migrationsB, migrationsA []*Migration) []*Migration {
	// Create a map for fast lookup
	lookup := make(map[int]bool)
	for _, migration := range migrationsA {
		lookup[migration.Number] = true
	}

	// Filter migrationsB
	filteredMigrations := make([]*Migration, 0)
	for _, migration := range migrationsB {
		if _, ok := lookup[migration.Number]; ok {
			filteredMigrations = append(filteredMigrations, migration)
		}
	}

	return filteredMigrations

}
