package migrator

func SelectDatabase(database string) Database {
	switch database {
	case "postgres", "postgresql":
		return PostgresDB
	case "mysql":
		return MySQLDB
	case "sqlite", "sqlite3":
		return SQLiteDB
	default:
		return SQLiteDB
	}
}

func SelectExtention(database Database) string {
	var ext string
	switch database {
	case PostgresDB:
		ext = "sql"
	case MySQLDB:
		ext = "sql"
	case SQLiteDB:
		ext = "sqlite"
	default:
		ext = "sql"
	}

	return ext
}
