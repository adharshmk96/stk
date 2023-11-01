package db

import (
	"database/sql"
	"sync"

	"github.com/adharshmk96/stktemplate/server/infra"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var (
	sqliteInstance *sql.DB
	sqliteOnce     sync.Once
)

// GetSqliteConnection returns a singleton database connection
func GetSqliteConnection() *sql.DB {
	filepath := viper.GetString(infra.ENV_SQLITE_FILEPATH)
	sqliteOnce.Do(func() {
		db, err := sql.Open("sqlite3", filepath)
		if err != nil {
			panic(err)
		}
		sqliteInstance = db
	})
	return sqliteInstance
}

// ResetSqliteConnection resets the singleton database connection
func ResetSqliteConnection() {
	sqliteInstance = nil
	sqliteOnce = sync.Once{}
}
