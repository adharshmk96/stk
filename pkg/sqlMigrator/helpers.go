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
func LoadUnappliedMigrations(ctx *Context) []*MigrationFileEntry {
	migrations := []*MigrationFileEntry{}

	for _, migration := range ctx.Migrations {
		if !migration.Committed {
			migrations = append(migrations, migration)
		}
	}

	return migrations
}

func LoadAppliedMigrations(ctx *Context) []*MigrationFileEntry {
	migrations := []*MigrationFileEntry{}

	for _, migration := range ctx.Migrations {
		if migration.Committed {
			migrations = append(migrations, migration)
		}
	}

	return migrations
}

func LastMigration(ctx *Context) *MigrationFileEntry {
	lastMigration := &MigrationFileEntry{}
	if len(ctx.Migrations) > 0 {
		lastMigration = ctx.Migrations[len(ctx.Migrations)-1]
	}

	return lastMigration
}
