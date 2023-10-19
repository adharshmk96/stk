package sqlmigrator

type MigrateDatabase interface {
	Exec(query string) error
}

// sqlite implementation
type SqliteMigrateDatabase struct {
	DatabasePath string
}

func (db *SqliteMigrateDatabase) Exec(query string) error {
	return nil
}
