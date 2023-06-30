package migrator

import "errors"

var (
	ErrInvalidFormat        = errors.New("invalid_format")
	ErrParsingMigrationType = errors.New("invalid_migration_type")

	ErrReadingFileNames      = errors.New("error_reading_file_names")
	ErrParsingMigrations     = errors.New("error_parsing_migrations")
	ErrCreatingMigrationFile = errors.New("error_creating_migration_file")

	ErrReadingLastAppliedMigration = errors.New("error_reading_last_applied_migration")
)
