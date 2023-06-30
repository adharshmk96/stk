package migrator

import "log"

type GeneratorConfig struct {
	RootDirectory string
	Database      string
	Name          string
	NumToGenerate int
	DryRun        bool
}

func Generate(config GeneratorConfig) error {
	// Select based on the database
	database := SelectDatabase(config.Database)
	log.Println("selected database: ", database)
	workDirectory := openDirectory(config.RootDirectory, database)
	log.Println("workdir: ", workDirectory)

	filePaths, err := getMigrationFilePathsByGroup(workDirectory, MigrationUp)
	if err != nil {
		return ErrReadingFileNames
	}

	lastMigrationNumber := 0

	if len(filePaths) > 0 {
		lastMigrationNumber, err = getLastMigrationNumber(filePaths)
		if err != nil {
			return err
		}

	}

	nextMigrations := generateNextMigrations(lastMigrationNumber, config.Name, config.NumToGenerate)

	for _, migration := range nextMigrations {
		fileName := migrationToFilename(migration) + "." + GetExtention(database)

		if config.DryRun {
			log.Println("dry run: ", fileName)
			continue
		}

		log.Println("generating file: ", fileName)
		err := createMigrationFile(workDirectory, fileName)
		if err != nil {
			return ErrCreatingMigrationFile
		}
	}

	return nil

}

func generateNextMigrations(lastNumber int, name string, total int) []*Migration {
	migrations := make([]*Migration, 0, total)
	for i := 0; i < total; i++ {
		migrations = append(migrations, &Migration{
			Number: lastNumber + i + 1,
			Name:   name,
			Type:   MigrationUp,
		})
		migrations = append(migrations, &Migration{
			Number: lastNumber + i + 1,
			Name:   name,
			Type:   MigrationDown,
		})
	}
	return migrations
}

func getLastMigrationNumber(filePaths []string) (int, error) {
	migrations, err := parseMigrationsFromFilePaths(filePaths)
	if err != nil {
		return 0, ErrParsingMigrations
	}

	sortMigrations(migrations)

	lastMigrationNumber := migrations[len(migrations)-1].Number
	return lastMigrationNumber, nil
}
