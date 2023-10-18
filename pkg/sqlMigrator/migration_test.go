package sqlmigrator_test

import (
	"testing"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/stretchr/testify/assert"
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
			assert.NoError(t, err)

			assert.Equal(t, c.expectedNumber, rawMigration.Number)

			assert.Equal(t, c.expectedName, rawMigration.Name)
		}
	})
}

func TestRawMigrationString(t *testing.T) {
	t.Run("outputs correct migration string", func(t *testing.T) {
		tc := []struct {
			rawMigration sqlmigrator.Migration
			expected     string
		}{
			{
				rawMigration: sqlmigrator.Migration{
					Number: 1,
					Name:   "create_users_table",
				},
				expected: "1_create_users_table",
			},
			{
				rawMigration: sqlmigrator.Migration{
					Number: 2,
					Name:   "",
				},
				expected: "2",
			},
		}

		for _, c := range tc {
			assert.Equal(t, c.expected, c.rawMigration.String())
		}
	})
}
