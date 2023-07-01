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
	FSRepo        FileRepo
}

func Generate(config GeneratorConfig) error {

	fsRepo := config.FSRepo

	// load migration from files
	migrations, err := fsRepo.LoadMigrationsFromFile(MigrationUp)
	if err != nil {
		log.Println("error loading migrations from file: ", err)
		return ErrLoadingMigrations
	}

	var lastMigrationNumber int

	if len(migrations) == 0 {
		lastMigrationNumber = 0
	} else {
		lastMigrationNumber = migrations[len(migrations)-1].Number
	}

	nextMigrations := GenerateNextMigrations(lastMigrationNumber, config.Name, config.NumToGenerate)

	if config.DryRun {
		dryRunGeneration(nextMigrations)
		return nil
	}

	GenerateMigrationFiles(fsRepo, nextMigrations)

	if config.Fill {
		FillMigrationFiles(fsRepo, config, nextMigrations)
	}

	return nil

}

func dryRunGeneration(migrations []*Migration) {
	for _, migration := range migrations {
		fileName := MigrationToFilename(migration)
		log.Println("dry run: ", fileName)
	}
}

func GenerateMigrationFiles(fsRepo FileRepo, migrations []*Migration) error {
	for _, migration := range migrations {
		err := fsRepo.CreateMigrationFile(migration)
		log.Println("generating file: ", migration.Path)
		if err != nil {
			return ErrCreatingMigrationFile
		}
	}
	return nil
}

func FillMigrationFiles(fsRepo FileRepo, config GeneratorConfig, migrations []*Migration) error {
	for _, migration := range migrations {
		var content string
		if migration.Type == MigrationUp {
			content = fmt.Sprintf("CREATE TABLE table_%d_%s (id INT PRIMARY KEY)", migration.Number, config.Name)
			migration.Query = content
		} else {
			content = fmt.Sprintf("DROP TABLE table_%d_%s", migration.Number, config.Name)
			migration.Query = content
		}
		err := fsRepo.WriteMigrationToFile(migration)
		if err != nil {
			return ErrCreatingMigrationFile
		}
	}
	return nil
}

func GenerateNextMigrations(lastNumber int, name string, total int) []*Migration {
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
