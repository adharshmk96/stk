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

type RawMigration struct {
	Number int
	Name   string
}

func ParseRawMigration(migration string) (*RawMigration, error) {
	parts := strings.Split(migration, "_")
	if len(parts) == 0 {
		return nil, ErrInvalidMigration
	}

	name := strings.Join(parts[1:], "_")

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, ErrInvalidMigration
	}

	rawMigration := &RawMigration{
		Name:   name,
		Number: number,
	}

	return rawMigration, nil
}

func (r *RawMigration) String() string {
	return fmt.Sprintf("%d_%s", r.Number, r.Name)
}
