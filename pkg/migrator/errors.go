package migrator

import "errors"

var (
	ErrInvalidFormat        = errors.New("invalid_format")
	ErrParsingMigrationType = errors.New("invalid_migration_type")

	ErrReadingFileNames      = errors.New("error_reading_file_names")
	ErrParsingMigrations     = errors.New("error_parsing_migrations")
	ErrCreatingMigrationFile = errors.New("error_creating_migration_file")

	ErrNoMigrationsToApply         = errors.New("error_no_migrations_to_apply")
	ErrReadingLastAppliedMigration = errors.New("error_reading_last_applied_migration")

	// Storage
	ErrMigrationTableDoesNotExist = errors.New("error_migration_table_does_not_exist")
	ErrDatabaseNotInitialized     = errors.New("error_database_not_initialized")
)
