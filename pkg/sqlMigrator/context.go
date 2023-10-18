package sqlmigrator

import (
	"path"

	"github.com/spf13/viper"
)

type MigrationType string

const (
	MigrationUp   MigrationType = "up"
	MigrationDown MigrationType = "down"
)

type Database string

const (
	PostgresDB Database = "postgres"
	MySQLDB    Database = "mysql"
	SQLiteDB   Database = "sqlite"
)

type Context struct {
	WorkDir  string
	Database Database
	DryRun   bool
}

func NewMigratorContext(dry bool) *Context {
	rootDirectory := viper.GetString("migrator.workdir")
	dbChoice := viper.GetString("migrator.database")

	dbType := SelectDatabase(dbChoice)
	subDir := SelectSubDirectory(dbType)

	workDir := path.Join(rootDirectory, subDir)

	return &Context{
		WorkDir:  workDir,
		Database: dbType,

		DryRun: dry,
	}
}
