package migrator

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MigrationType string

const (
	MigrationUp   MigrationType = "up"
	MigrationDown MigrationType = "down"
)

type Database string

const (
	PostgresDB Database = "postgres"
	MySQLDB    Database = "mysql"
	SQLiteDB   Database = "sqlite"
)

type Migration struct {
	Number  int
	Name    string
	Type    MigrationType
	Path    string
	Query   string
	Created time.Time
}

func parseMigrationType(s string) (MigrationType, error) {
	switch s {
	case "up":
		return MigrationUp, nil
	case "down":
		return MigrationDown, nil
	default:
		return "", ErrParsingMigrationType
	}
}

func ParseMigrationFromString(s string) (*Migration, error) {
	parts := strings.Split(s, "_")
	if len(parts) < 2 {
		return nil, ErrInvalidFormat
	}

	// parse 0th index for number
	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, ErrInvalidFormat
	}

	// remove 0th index from parts
	parts = parts[1:]

	// parse last part for migration type
	mType := parts[len(parts)-1]
	migrationType, err := parseMigrationType(mType)
	if err != nil {
		return nil, ErrInvalidFormat
	}

	// parse the remaining part for migration name
	name := ""
	if len(parts) > 1 {
		name = strings.Join(parts[:len(parts)-1], "_")
	}

	return &Migration{
		Number: number,
		Name:   name,
		Type:   migrationType,
	}, nil
}

func sortMigrations(migrations []*Migration) {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Number < migrations[j].Number
	})
}

func MigrationToFilename(migration *Migration) string {
	migration.Name = strings.ReplaceAll(migration.Name, " ", "_")
	return fmt.Sprintf("%06d_%s_%s", migration.Number, migration.Name, migration.Type)
}

type DatabaseRepo interface {
	// Load last applied migration entry from the migration table
	// - Creates migration table if not exists
	// - Returns nil if no entry found
	LoadLastAppliedMigration() (*Migration, error)
	// Load all migration entries from the migration table
	// - Creates migration table if not exists
	LoadMigrations() ([]*Migration, error)
	// Create a migration table if not exists
	CreateMigrationTableIfNotExists() error
	// Get the last applied migration from the migration table
	GetLastAppliedMigration() (*Migration, error)
	// Apply a migration to the database and add an entry to the migration table
	ApplyMigration(migration *Migration) error
	// Delete the migration table
	DeleteMigrationTable() error
}

type FileRepo interface {
	// Create a migration directory if not exists
	OpenDirectory() error
	// Load migrations from directoryy
	LoadMigrationsFromFile(migrationType MigrationType) ([]*Migration, error)
	// Read all the files from the migration directory
	GetMigrationFilePathsByType(migrationType MigrationType) ([]string, error)
	// Create a migration file
	CreateMigrationFile(migration *Migration) error
	// Write to File
	WriteMigrationToFile(migration *Migration) error
	// Read Query from File
	LoadMigrationQuery(migration *Migration) error
}
