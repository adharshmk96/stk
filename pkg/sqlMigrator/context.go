package sqlmigrator

import (
	"path"

	"github.com/spf13/viper"
)

const (
	DEFAULE_LOG_FILE = "migrations.log"
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
	LogFile  string
	Database Database
	DryRun   bool
}

func NewMigratorContext(dry bool) *Context {
	rootDirectory := viper.GetString("migrator.workdir")
	dbChoice := viper.GetString("migrator.database")
	logFile := getFirst(viper.GetString("migrator.logfile"), DEFAULE_LOG_FILE)

	dbType := SelectDatabase(dbChoice)
	subDir := SelectSubDirectory(dbType)

	workDir := path.Join(rootDirectory, subDir)

	ctx := &Context{
		WorkDir:  workDir,
		Database: dbType,
		LogFile:  logFile,
		DryRun:   dry,
	}

	if !dry {
		err := InitializeMigrationsLog(ctx)
		if err != nil {
			return nil
		}
	}

	return ctx
}
