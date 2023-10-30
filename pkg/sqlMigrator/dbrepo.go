package sqlmigrator

type DBRepo interface {
	Exec(query string) error
	PushHistory(migration *MigrationDBEntry) error
	LoadHistory() ([]*MigrationDBEntry, error)
	InitMigrationTable() error
	DeleteMigrationTable() error
}
