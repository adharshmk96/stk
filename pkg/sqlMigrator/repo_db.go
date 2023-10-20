package sqlmigrator

type DBRepo interface {
	Exec(query string) error
	LoadHistory() ([]*MigrationEntry, error)
}

// sqlite implementation
type SqliteMigrateDatabase struct {
	DatabasePath string
}

func NewSQLiteRepo(databasePath string) *SqliteMigrateDatabase {
	return &SqliteMigrateDatabase{
		DatabasePath: databasePath,
	}
}

func (db *SqliteMigrateDatabase) Exec(query string) error {
	return nil
}

func (db *SqliteMigrateDatabase) LoadHistory() ([]*MigrationEntry, error) {
	return []*MigrationEntry{}, nil
}
