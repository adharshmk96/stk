package sqlmigrator

import (
	"os"
	"path"
)

func InitializeMigrationsLog(ctx *Context) error {
	err := os.MkdirAll(ctx.WorkDir, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := path.Join(ctx.WorkDir, ctx.LogFile)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			// create file
			_, err := os.Create(filePath)
			if err != nil {
				return err
			}
		} else {
			// other error
			return err
		}
	}

	return nil
}

func loadLastMigrationFromLog(ctx *Context) (*Migration, error) {
	filePath := path.Join(ctx.WorkDir, ctx.LogFile)
	lastLine, err := readLastLine(filePath)
	if err != nil {
		return nil, err
	}
	if lastLine == "" {
		return &Migration{}, nil
	}

	lastMigration, err := ParseRawMigration(lastLine)
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
