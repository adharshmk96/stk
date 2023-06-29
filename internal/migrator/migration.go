package migrator

import (
	"strconv"
	"strings"
)

type MigrationType string

const (
	MigrationUp   MigrationType = "up"
	MigrationDown MigrationType = "down"
)

type Migration struct {
	Number int
	Name   string
	Type   MigrationType
}

func ParseMigrationType(s string) (MigrationType, error) {
	switch s {
	case "up":
		return MigrationUp, nil
	case "down":
		return MigrationDown, nil
	default:
		return "", ErrParsingMigrationType
	}
}

func ParseMigration(s string) (*Migration, error) {
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
	migrationType, err := ParseMigrationType(mType)
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

func ParseMigrationsFromFilenames(filenames []string) ([]*Migration, error) {
	migrations := make([]*Migration, 0, len(filenames))
	for _, filename := range filenames {
		migration, err := ParseMigration(filename)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}
