package dbrepo

import "database/sql"

func setupDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	return db
}
