package dbrepo

import (
	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/spf13/viper"
)

const (
	MIGRATION_TABLE_NAME = "stk_migrations"
)

func SelectDBRepo(database sqlmigrator.Database) sqlmigrator.DBRepo {
	switch database {
	case sqlmigrator.SQLiteDB:
		viper.SetDefault("migrator.database.filepath", "migrations.db")
		filePath := viper.GetString("migrator.database.filepath")
		return NewSQLiteRepo(filePath)
	default:
		viper.SetDefault("migrator.database.filepath", "migrations.db")
		filePath := viper.GetString("migrator.database.filepath")
		return NewSQLiteRepo(filePath)
	}
}
