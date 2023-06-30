package migrator

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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

func SelectDatabase(database string) Database {
	switch database {
	case "postgres":
		return PostgresDB
	case "mysql":
		return MySQLDB
	case "sqlite":
		return SQLiteDB
	default:
		return SQLiteDB
	}
}

func GetExtention(database Database) string {
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

type Migration struct {
	Number int
	Name   string
	Type   MigrationType
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

func parseMigration(s string) (*Migration, error) {
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

func parseMigrationsFromFilenames(filenames []string) ([]*Migration, error) {
	migrations := make([]*Migration, 0, len(filenames))
	for _, filename := range filenames {
		migration, err := parseMigration(filename)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func sortMigrations(migrations []*Migration) {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Number < migrations[j].Number
	})
}

func generateNextMigrations(lastNumber int, name string, total int) []*Migration {
	migrations := make([]*Migration, 0, total)
	for i := 0; i < total; i++ {
		migrations = append(migrations, &Migration{
			Number: lastNumber + i + 1,
			Name:   name,
			Type:   MigrationUp,
		})
		migrations = append(migrations, &Migration{
			Number: lastNumber + i + 1,
			Name:   name,
			Type:   MigrationDown,
		})
	}
	return migrations
}

func migrationToFilename(migration *Migration) string {
	migration.Name = strings.ReplaceAll(migration.Name, " ", "_")
	return fmt.Sprintf("%06d_%s_%s", migration.Number, migration.Name, migration.Type)
}
