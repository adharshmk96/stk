package sqlmigrator_test

import (
	"testing"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
)

func TestParseRawMigration(t *testing.T) {
	t.Run("parses a correct raw migration", func(t *testing.T) {
		tc := []struct {
			rawMigrationString string
			expectedNumber     int
			expectedName       string
		}{
			{
				rawMigrationString: "1_create_users_table",
				expectedNumber:     1,
				expectedName:       "create_users_table",
			},
			{
				rawMigrationString: "2",
				expectedNumber:     2,
				expectedName:       "",
			},
		}

		for _, c := range tc {
			rawMigration, err := sqlmigrator.ParseRawMigration(c.rawMigrationString)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rawMigration.Number != c.expectedNumber {
				t.Fatalf("expected number %d, got %d", c.expectedNumber, rawMigration.Number)
			}

			if rawMigration.Name != c.expectedName {
				t.Fatalf("expected name %s, got %s", c.expectedName, rawMigration.Name)
			}
		}
	})
}
