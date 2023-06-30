package migrator

func SelectSubDirectory(database Database) string {
	var subDirectory string
	switch database {
	case PostgresDB:
		subDirectory = "postgres"
	case MySQLDB:
		subDirectory = "mysql"
	case SQLiteDB:
		subDirectory = "sqlite"
	default:
		subDirectory = "sqlite"
	}

	return subDirectory
}
