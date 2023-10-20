package sqlmigrator

import (
	"fmt"
	"os"
	"path"
)

type Generator struct {
	Name          string
	NumToGenerate int
	DryRun        bool
	Fill          bool
}

func NewGenerator(name string, numToGenerate int, fill bool) *Generator {
	return &Generator{
		Name:          name,
		NumToGenerate: numToGenerate,
		Fill:          fill,
	}
}

func (g *Generator) Generate(ctx *Context) ([]string, error) {
	generatedFiles := []string{}
	// Assumes that the log file exists, It is generated when context is initialized
	lastMigration := LastMigration(ctx)

	nextMigrations := GenerateNextMigrations(lastMigration.Number, g.Name, g.NumToGenerate)
	if ctx.DryRun {
		dryRunGeneration(nextMigrations)
		return generatedFiles, nil
	}

	for _, migration := range nextMigrations {
		migString := migration.String()
		extention := SelectExtention(ctx.Database)

		upFileName, downFileName := migration.FileNames(extention)

		migration.UpFilePath = path.Join(ctx.WorkDir, upFileName)
		upFileContent := ""
		if g.Fill {
			upFileContent = fmt.Sprintf("CREATE TABLE sample_%s_table (id INTEGER PRIMARY KEY AUTOINCREMENT, created DATETIME DEFAULT CURRENT_TIMESTAMP);", migString)
		}
		migration.DownFilePath = path.Join(ctx.WorkDir, downFileName)
		downFileContent := ""
		if g.Fill {
			downFileContent = fmt.Sprintf("DROP TABLE sample_%s_table;", migString)
		}
		err := CreateFile(migration.UpFilePath, upFileContent)
		if err != nil {
			return generatedFiles, err
		}

		err = CreateFile(migration.DownFilePath, downFileContent)
		if err != nil {
			return generatedFiles, err
		}

		ctx.Migrations = append(ctx.Migrations, migration)
		generatedFiles = append(generatedFiles, migration.UpFilePath, migration.DownFilePath)
	}

	return generatedFiles, nil
}

func (g *Generator) Clean(ctx *Context) ([]string, error) {
	removedFiles := []string{}

	uncommitedMigrations, err := LoadUncommitedMigrations(ctx)
	if err != nil {
		return nil, err
	}
	if ctx.DryRun {
		dryRunGeneration(uncommitedMigrations)
		return removedFiles, nil
	}

	for _, migration := range uncommitedMigrations {
		upFileName, downFileName := migration.FileNames(SelectExtention(ctx.Database))

		upFilePath := path.Join(ctx.WorkDir, upFileName)
		downFilePath := path.Join(ctx.WorkDir, downFileName)

		err := os.Remove(upFilePath)
		if err != nil {
			return removedFiles, err
		}

		err = os.Remove(downFilePath)
		if err != nil {
			return removedFiles, err
		}

		removedFiles = append(removedFiles, upFilePath, downFilePath)
	}

	ctx.Migrations = ctx.Migrations[:len(ctx.Migrations)-len(uncommitedMigrations)]
	return removedFiles, nil
}

func dryRunGeneration(migrations []*MigrationFileEntry) {
	for _, migration := range migrations {
		fileName := migration.EntryString()
		fmt.Println("up\t:", fileName+"_up.sql")
		fmt.Println("down:\t:", fileName+"_down.sql")
	}
}

func GenerateNextMigrations(lastMigrationNumber int, name string, numToGenerate int) []*MigrationFileEntry {
	var nextMigrations []*MigrationFileEntry

	startNumber := lastMigrationNumber + 1
	endNumber := lastMigrationNumber + numToGenerate

	for i := startNumber; i <= endNumber; i++ {
		nextMigration := &MigrationFileEntry{
			Number: i,
			Name:   name,
		}

		nextMigrations = append(nextMigrations, nextMigration)
	}

	return nextMigrations
}
