package migrator

type GeneratorConfig struct {
	RootDirectory string
	Database      string
	Name          string
	NumToGenerate int
}

func Generate(config GeneratorConfig) error {
	// Select based on the database
	database := SelectDatabase(config.Database)
	subDirectory := OpenDirectory(database)

	fileNames, err := GetMigrationFileGroup(subDirectory, MigrationUp)
	if err != nil {
		return ErrReadingFileNames
	}

	migrations, err := ParseMigrationsFromFilenames(fileNames)
	if err != nil {
		return ErrParsingMigrations
	}

	SortMigrations(migrations)

	lastMigrationNumber := migrations[len(migrations)-1].Number

	nextMigrations := GenerateNextMigrations(lastMigrationNumber, config.Name, config.NumToGenerate)

	for _, migration := range nextMigrations {
		fileName := MigrationToFilename(migration) + "." + GetExtention(database)
		err := CreateMigrationFile(subDirectory, fileName)
		if err != nil {
			return ErrCreatingMigrationFile
		}
	}

	return nil

}
