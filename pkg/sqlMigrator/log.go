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

// TODO: optimize to read from last, and break when commit status is true
func LoadUncommitedMigrations(ctx *Context) ([]*MigrationEntry, error) {
	migrations := []*MigrationEntry{}

	for _, migration := range ctx.Migrations {
		if !migration.Committed {
			migrations = append(migrations, migration)
		}
	}

	return migrations, nil
}
