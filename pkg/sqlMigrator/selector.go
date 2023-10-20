package sqlmigrator

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

func SelectSubDirectory(database Database) string {
	var subDir string
	switch database {
	case PostgresDB:
		subDir = "postgres"
	case MySQLDB:
		subDir = "mysql"
	case SQLiteDB:
		subDir = "sqlite"
	default:
		subDir = "sqlite"
	}

	return subDir
}

func SelectDBRepo(database Database, path string) DBRepo {
	var repo DBRepo
	switch database {
	case SQLiteDB:
		repo = NewSQLiteRepo(path)
	default:
		repo = NewSQLiteRepo(path)
	}

	return repo
}
