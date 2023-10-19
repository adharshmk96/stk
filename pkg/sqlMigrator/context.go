package sqlmigrator

import (
	"path"

	"github.com/spf13/viper"
)

const (
	DEFAULT_LOG_FILE = ".commit-status"
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

func DefaultContextConfig() (string, Database, string) {
	rootDirectory := viper.GetString("migrator.workdir")
	dbChoice := viper.GetString("migrator.database")
	logFile := getFirst(viper.GetString("migrator.logfile"), DEFAULT_LOG_FILE)

	dbType := SelectDatabase(dbChoice)
	subDir := SelectSubDirectory(dbType)

	workDir := path.Join(rootDirectory, subDir)

	return workDir, dbType, logFile
}

func NewMigratorContext(workDir string, dbType Database, logFile string, dry bool) *Context {

	ctx := &Context{
		WorkDir:  workDir,
		Database: dbType,
		LogFile:  logFile,
		DryRun:   dry,
	}

	err := InitializeMigrationsFolder(ctx)
	if err != nil {
		return nil
	}

	return ctx
}
