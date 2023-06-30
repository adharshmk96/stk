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
	// 2. Read all migrations from file system
	workDirectory := openDirectory(config.RootDirectory, database)
	log.Println("workdir: ", workDirectory)

	filePaths, err := getMigrationFilePathsByGroup(workDirectory, MigrationUp)
	if err != nil {
		return ErrReadingFileNames
	}

	log.Println("filenames: ", filePaths)

	// 3. Sort the migrations
	// 4. Find the next migrations to apply
	// 5. Read migration queries from files
	// 6. Apply migrations and add entries to database

	return nil

}
