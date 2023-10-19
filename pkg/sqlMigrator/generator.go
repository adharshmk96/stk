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

func NewGenerator(name string, numToGenerate int, dryRun bool, fill bool) *Generator {
	return &Generator{
		Name:          name,
		NumToGenerate: numToGenerate,
		Fill:          fill,
	}
}

func (g *Generator) Generate(ctx *Context) error {
	// Assumes that the log file exists, It is generated when context is initialized
	lastMigration, err := loadLastMigrationFromLog(ctx)
	if err != nil {
		return err
	}

	nextMigrations := GenerateNextMigrations(lastMigration.Number, g.Name, g.NumToGenerate)
	if ctx.DryRun {
		dryRunGeneration(nextMigrations)
		return nil
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
			return err
		}

		err = createFile(downFilePath, downFileContent)
		if err != nil {
			return err
		}

		err = writeMigrationToLog(ctx, migrationName)
		if err != nil {
			return err
		}
	}

	return nil
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
