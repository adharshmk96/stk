package migrator

import "log"

type MigratorConfig struct {
	RootDirectory string
	Database      string
	NumToMigrate  int
	DryRun        bool
	dbRepo        DatabaseRepo
}

func MigrateUp(config *MigratorConfig) (err error) {
	// Select based on the database
	database := SelectDatabase(config.Database)
	log.Println("selected database: ", database)

	// 1. Read last applied migration from database
	lastAppliedMigration, err := config.dbRepo.GetLastAppliedMigration()
	if err != nil {
		return ErrReadingLastAppliedMigration
	}

	// 2. Read all migrations from file system
	workDirectory := openDirectory(config.RootDirectory, database)
	log.Println("workdir: ", workDirectory)

	filePaths, err := getMigrationFilePathsByGroup(workDirectory, MigrationUp)
	if err != nil {
		return ErrReadingFileNames
	}

	log.Println("filenames: ", filePaths)

	// 3. Parse migrations from file paths
	migrations, err := parseMigrationsFromFilePaths(filePaths)
	if err != nil {
		return ErrParsingMigrations
	}

	// 3. Sort the migrations
	sortMigrations(migrations)

	log.Println("migrations: ", migrations)
	log.Println("last : ", lastAppliedMigration)

	// // 4. Find the next migrations to apply
	// nextMigrations := findUpMigrationsToApply(lastAppliedMigration, migrations, config.NumToMigrate)
	// // 5. Read migration queries from files
	// queries := readMigrationQueries(nextMigrations)
	// // 6. Apply migrations and add entries to database
	// apply all as one transaction. if one fails, rollback all
	// err = applyMigrations(config.dbRepo, nextMigrations, queries, config.DryRun)

	return nil

}

func findNextMigrationIndex(migrations []*Migration, number int) int {
	for i, v := range migrations {
		if v.Number > number {
			return i
		}
	}
	return -1
}

func findUpMigrationsToApply(lastAppliedMigration *MigrationEntry, migrations []*Migration, numberToMigrate int) []*Migration {
	idx := findNextMigrationIndex(migrations, lastAppliedMigration.Number)
	return migrations[idx : idx+numberToMigrate]
}
