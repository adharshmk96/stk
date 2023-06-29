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
	workDirectory := OpenDirectory(config.RootDirectory, database)
	log.Println("workdir: ", workDirectory)

	fileNames, err := GetMigrationFileGroup(workDirectory, MigrationUp)
	if err != nil {
		return ErrReadingFileNames
	}

	lastMigrationNumber := 0

	if len(fileNames) > 0 {
		lastMigrationNumber, err = getLastMigrationNumber(fileNames)
		if err != nil {
			return err
		}

	}

	nextMigrations := GenerateNextMigrations(lastMigrationNumber, config.Name, config.NumToGenerate)

	for _, migration := range nextMigrations {
		fileName := MigrationToFilename(migration) + "." + GetExtention(database)

		if config.DryRun {
			log.Println("dry run: ", fileName)
			continue
		}

		log.Println("generating file: ", fileName)
		err := CreateMigrationFile(workDirectory, fileName)
		if err != nil {
			return ErrCreatingMigrationFile
		}
	}

	return nil

}

func getLastMigrationNumber(fileNames []string) (int, error) {
	migrations, err := ParseMigrationsFromFilenames(fileNames)
	if err != nil {
		return 0, ErrParsingMigrations
	}

	SortMigrations(migrations)

	lastMigrationNumber := migrations[len(migrations)-1].Number
	return lastMigrationNumber, nil
}
