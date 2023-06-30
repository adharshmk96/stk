package database

const (
	sqlMigrationTable = "migdb_migration"

	sqliteMigrationTableExists      = "SELECT 1 FROM sqlite_master WHERE type='table' AND name=?;"
	sqliteSelectMigrationEntries    = "SELECT name, migtype, created FROM " + sqlMigrationTable + " ORDER BY id ASC"
	sqliteDropMigrationTable        = "DROP TABLE IF EXISTS " + sqlMigrationTable
	sqliteLastAppliedMigrationEntry = "SELECT name, migtype, created FROM " + sqlMigrationTable + " ORDER BY id DESC LIMIT 1"
	sqliteInsertMigrationEntry      = "INSERT INTO " + sqlMigrationTable + " (name, migtype) VALUES ($1, $2)"
	sqliteMigrationSchema           = `CREATE TABLE IF NOT EXISTS migdb_migration (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		migtype VARCHAR(5) NOT NULL,
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
)
