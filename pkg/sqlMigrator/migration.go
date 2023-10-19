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

type MigrationEntry struct {
	Number       int
	Name         string
	CommitStatus bool
	Up           string
	Down         string
}

func ParseMigrationEntry(migrationEntry string) (*MigrationEntry, error) {
	parts := strings.Split(migrationEntry, "_")
	partLength := len(parts)

	if partLength == 0 {
		return nil, ErrInvalidMigration
	}

	commit_status := parts[partLength-1]
	if commit_status != "up" && commit_status != "down" {
		return nil, ErrInvalidMigration
	}

	name := strings.Join(parts[1:partLength-1], "_")

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, ErrInvalidMigration
	}

	rawMigration := &MigrationEntry{
		Name:         name,
		Number:       number,
		CommitStatus: commit_status == "up",
	}

	return rawMigration, nil
}

func (r *MigrationEntry) String() string {
	entryString := fmt.Sprintf("%d", r.Number)
	if r.Name != "" {
		entryString += "_" + r.Name
	}
	if r.CommitStatus {
		entryString += "_up"
	} else {
		entryString += "_down"
	}
	return entryString
}

func (r *MigrationEntry) FileNames(extention string) (string, string) {
	fileName := fmt.Sprintf("%d", r.Number)
	if r.Name != "" {
		fileName += "_" + r.Name
	}
	upFileName := fileName + "_up." + extention
	downFileName := fileName + "_down." + extention
	return upFileName, downFileName
}
