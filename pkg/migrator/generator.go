package migrator

import (
	"fmt"
	"log"
)

type GeneratorConfig struct {
	Name          string
	NumToGenerate int
	DryRun        bool
	Fill          bool
	Database      Database
	FSRepo        FileRepo
}

func Generate(config GeneratorConfig) error {

	fsRepo := config.FSRepo

	// load migration from files
	migrations, err := fsRepo.LoadMigrationsFromFile(MigrationUp)
	if err != nil {
		log.Fatalln("error loading migrations from file: ", err)
		return ErrLoadingMigrations
	}

	sortMigrations(migrations)

	var lastMigrationNumber int

	if len(migrations) == 0 {
		lastMigrationNumber = 0
	} else {
		lastMigrationNumber = migrations[len(migrations)-1].Number
	}

	nextMigrations := generateNextMigrations(lastMigrationNumber, config.Name, config.NumToGenerate)

	if config.DryRun {
		dryRunGeneration(nextMigrations)
		return nil
	}

	generateMigrationFiles(fsRepo, nextMigrations)

	if config.Fill {
		fillMigrationFiles(fsRepo, config, nextMigrations)
	}

	return nil

}

func dryRunGeneration(migrations []*Migration) {
	for _, migration := range migrations {
		fileName := MigrationToFilename(migration)
		log.Println("dry run: ", fileName)
	}
}

func generateMigrationFiles(fsRepo FileRepo, migrations []*Migration) error {
	for _, migration := range migrations {
		err := fsRepo.CreateMigrationFile(migration)
		log.Println("generating file: ", migration.Path)
		if err != nil {
			return ErrCreatingMigrationFile
		}
	}
	return nil
}

func fillMigrationFiles(fsRepo FileRepo, config GeneratorConfig, migrations []*Migration) error {
	for _, migration := range migrations {
		var content string
		if migration.Type == MigrationUp {
			content = fmt.Sprintf("CREATE TABLE %d_%s_%s (id INT PRIMARY KEY)", migration.Number, config.Name, migration.Type)
			migration.Query = content
		} else {
			content = fmt.Sprintf("DROP TABLE %d_%s_%s", migration.Number, config.Name, migration.Type)
			migration.Query = content
		}
		err := fsRepo.WriteMigrationToFile(migration)
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
