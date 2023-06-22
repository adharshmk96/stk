package db

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// Define Singleton sqlite_instance
var sqlite_instance *sql.DB
var sqliteOnce sync.Once

// DB returns a singleton database connection
func GetSqliteConnection(filepath string) *sql.DB {
	sqliteOnce.Do(func() {
		db, err := sql.Open("sqlite3", filepath)
		if err != nil {
			panic(err)
		}
		sqlite_instance = db
	})
	return sqlite_instance
}
