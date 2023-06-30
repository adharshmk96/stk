package migrator

import "log"

type MigratorConfig struct {
	RootDirectory string
	Database      string
	NumToMigrate  int
	DryRun        bool
	DBRepo        DatabaseRepo
}

func MigrateUp(config *MigratorConfig) (err error) {
	// Select based on the database
	database := SelectDatabase(config.Database)
	log.Println("selected database: ", database)

	if config.DBRepo == nil {
		log.Fatalf("database is not initialized")
		return ErrDatabaseNotInitialized
	}
	// Read last applied migration from database
	err = config.DBRepo.CreateMigrationTableIfNotExists()
	if err != nil {
		return ErrMigrationTableDoesNotExist
	}
	lastAppliedMigration, err := config.DBRepo.GetLastAppliedMigration()
	if err != nil {
		return ErrReadingLastAppliedMigration
	}

	// Read all migrations from file system
	workDirectory := openDirectory(config.RootDirectory, database)
	log.Println("workdir: ", workDirectory)

	filePaths, err := getMigrationFilePathsByGroup(workDirectory, MigrationUp)
	// if len(filePaths) == 0 {
	// 	log.Println("no migrations to apply...")
	// 	return nil
	// }
	if err != nil {
		return ErrReadingFileNames
	}

	log.Println("filenames: ", filePaths)

	// Parse migrations from file paths
	migrations, err := parseMigrationsFromFilePaths(filePaths)
	if err != nil {
		return ErrParsingMigrations
	}

	// Sort the migrations
	sortMigrations(migrations)

	log.Println("migrations: ", migrations)
	log.Println("last : ", lastAppliedMigration)

	// Find the next migrations to apply
	migrationsToApply := findUpMigrationsToApply(lastAppliedMigration, migrations, config.NumToMigrate)
	// Read migration queries from files
	readMigrationQueries(migrationsToApply)

	if config.DryRun {
		log.Println("dry run: ")
		for _, migration := range migrationsToApply {
			log.Println(migration.Name)
		}
		return nil
	}

	// Apply migrations and add entries to database
	// TODO: apply all as one transaction. if one fails, rollback all
	for _, migration := range migrationsToApply {
		err = config.DBRepo.ApplyMigration(migration)
		if err != nil {
			log.Fatalln("error applying migration: ", migration.Number, err)
		}
	}

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

func findUpMigrationsToApply(lastMigration *Migration, migrations []*Migration, numberToMigrate int) []*Migration {
	if lastMigration == nil {
		return migrations
	}
	idx := findNextMigrationIndex(migrations, lastMigration.Number)
	if lastMigration.Type == MigrationDown {
		startIdx := max(idx-1, 0)
		lastIdx := min(startIdx+numberToMigrate, len(migrations))
		return migrations[startIdx:lastIdx]
	} else {
		startIdx := max(idx, 0)
		lastIdx := min(startIdx+numberToMigrate, len(migrations))
		return migrations[startIdx:lastIdx]

	}
}

func readMigrationQuery(migration *Migration) string {
	filePath := migration.Path
	query, err := readFileContents(filePath)
	if err != nil {
		log.Fatalln("error reading file contents: ", err)
		return ""
	}
	return query
}

func readMigrationQueries(migrations []*Migration) {
	for _, migration := range migrations {
		migration.Query = readMigrationQuery(migration)
	}
}
