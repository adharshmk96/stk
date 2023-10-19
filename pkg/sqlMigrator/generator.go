package sqlmigrator

import (
	"fmt"
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
	lastMigration, err := loadLastMigrationFromLog(ctx)
	if err != nil {
		return nil, err
	}

	nextMigrations := GenerateNextMigrations(lastMigration.Number, g.Name, g.NumToGenerate)
	if ctx.DryRun {
		dryRunGeneration(nextMigrations)
		return generatedFiles, nil
	}

	for _, migration := range nextMigrations {
		migrationName := migration.String()
		extention := SelectExtention(ctx.Database)
		upFile := migrationName + "_up." + extention
		downFile := migrationName + "_down." + extention

		upFilePath := path.Join(ctx.WorkDir, upFile)
		upFileContent := ""
		if g.Fill {
			upFileContent = fmt.Sprintf("CREATE TABLE sample_%s_table;", migrationName)
		}
		downFilePath := path.Join(ctx.WorkDir, downFile)
		downFileContent := ""
		if g.Fill {
			downFileContent = fmt.Sprintf("DROP TABLE sample_%s_table;", migrationName)
		}
		err := createFile(upFilePath, upFileContent)
		if err != nil {
			return generatedFiles, err
		}

		err = createFile(downFilePath, downFileContent)
		if err != nil {
			return generatedFiles, err
		}

		err = writeMigrationToLog(ctx, migrationName)
		if err != nil {
			return generatedFiles, err
		}

		generatedFiles = append(generatedFiles, upFilePath, downFilePath)
	}

	return generatedFiles, nil
}

func dryRunGeneration(migrations []*Migration) {
	for _, migration := range migrations {
		fileName := migration.String()
		fmt.Println("up\t:", fileName+"_up.sql")
		fmt.Println("down:\t:", fileName+"_down.sql")
	}
}

func GenerateNextMigrations(lastMigrationNumber int, name string, numToGenerate int) []*Migration {
	var nextMigrations []*Migration

	startNumber := lastMigrationNumber + 1
	endNumber := lastMigrationNumber + numToGenerate

	for i := startNumber; i <= endNumber; i++ {
		nextMigration := &Migration{
			Number: i,
			Name:   name,
		}

		nextMigrations = append(nextMigrations, nextMigration)
	}

	return nextMigrations
}
