package migrator

import "errors"

var (
	ErrInvalidFormat        = errors.New("invalid format")
	ErrParsingMigrationType = errors.New("invalid_migration_type")
)
