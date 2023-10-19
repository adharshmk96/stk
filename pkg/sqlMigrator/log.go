package sqlmigrator

import (
	"os"
	"path"
)

func InitializeMigrationsFolder(ctx *Context) error {
	err := os.MkdirAll(ctx.WorkDir, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := path.Join(ctx.WorkDir, ctx.LogFile)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			// create file
			file, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
		} else {
			// other error
			return err
		}
	}

	return nil
}

func loadLastMigrationFromLog(ctx *Context) (*MigrationEntry, error) {
	filePath := path.Join(ctx.WorkDir, ctx.LogFile)
	lastLine, err := readLastLine(filePath)
	if err != nil {
		return nil, err
	}
	if lastLine == "" {
		return &MigrationEntry{}, nil
	}

	lastMigration, err := ParseMigrationEntry(lastLine)
	if err != nil {
		return nil, err
	}

	return lastMigration, nil
}

func writeMigrationToLog(ctx *Context, migration string) error {
	filePath := path.Join(ctx.WorkDir, ctx.LogFile)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(migration + "\n")
	if err != nil {
		return err
	}

	return nil
}

// TODO: optimize to read from last, and break when commit status is true
func LoadUncommitedMigrationsFromLog(ctx *Context) ([]*MigrationEntry, error) {
	readLines, err := readLines(path.Join(ctx.WorkDir, ctx.LogFile))
	if err != nil {
		return nil, err
	}

	migrations := []*MigrationEntry{}
	for _, line := range readLines {
		migration, err := ParseMigrationEntry(line)
		if err != nil {
			return nil, err
		}

		if !migration.Committed {
			migrations = append(migrations, migration)
		}
	}

	return migrations, nil

}
