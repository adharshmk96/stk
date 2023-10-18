package sqlmigrator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidMigration = errors.New("invalid migration")
)

type Migration struct {
	Number int
	Name   string
	Up     string
	Down   string
}

func ParseRawMigration(migration string) (*Migration, error) {
	parts := strings.Split(migration, "_")
	if len(parts) == 0 {
		return nil, ErrInvalidMigration
	}

	name := strings.Join(parts[1:], "_")

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, ErrInvalidMigration
	}

	rawMigration := &Migration{
		Name:   name,
		Number: number,
	}

	return rawMigration, nil
}

func (r *Migration) String() string {
	if r.Name == "" {
		return fmt.Sprintf("%d", r.Number)
	}
	return fmt.Sprintf("%d_%s", r.Number, r.Name)
}
