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
			commitStatus       bool
		}{
			{
				rawMigrationString: "1_create_users_table_up",
				expectedNumber:     1,
				expectedName:       "create_users_table",
				commitStatus:       true,
			},
			{
				rawMigrationString: "2_up",
				expectedNumber:     2,
				expectedName:       "",
				commitStatus:       true,
			},
			{
				rawMigrationString: "3_create_posts_table_down",
				expectedNumber:     3,
				expectedName:       "create_posts_table",
				commitStatus:       false,
			},
			{
				rawMigrationString: "4_down",
				expectedNumber:     4,
				expectedName:       "",
				commitStatus:       false,
			},
		}

		for _, c := range tc {
			rawMigration, err := sqlmigrator.ParseMigrationEntry(c.rawMigrationString)
			assert.NoError(t, err)
			assert.Equal(t, c.expectedNumber, rawMigration.Number)
			assert.Equal(t, c.expectedName, rawMigration.Name)
			assert.Equal(t, c.commitStatus, rawMigration.CommitStatus)
		}
	})

	t.Run("returns an error if raw migration is invalid", func(t *testing.T) {
		tc := []struct {
			rawMigrationString string
		}{
			{
				rawMigrationString: "create_users_table_up",
			},
			{
				rawMigrationString: "1_create_users_table",
			},
			{
				rawMigrationString: "1",
			},
			{
				rawMigrationString: "1create_users_table",
			},
			{
				rawMigrationString: "create_users_table",
			},
			{
				rawMigrationString: "create_users_table1",
			},
			{
				rawMigrationString: "",
			},
			{
				rawMigrationString: "nameonly",
			},
		}

		for _, c := range tc {
			_, err := sqlmigrator.ParseMigrationEntry(c.rawMigrationString)
			assert.Error(t, err)
		}
	})
}

func TestRawMigrationString(t *testing.T) {
	t.Run("outputs correct migration string", func(t *testing.T) {
		tc := []struct {
			rawMigration sqlmigrator.MigrationEntry
			expected     string
		}{
			{
				rawMigration: sqlmigrator.MigrationEntry{
					Number:       1,
					Name:         "create_users_table",
					CommitStatus: true,
				},
				expected: "1_create_users_table_up",
			},
			{
				rawMigration: sqlmigrator.MigrationEntry{
					Number:       2,
					Name:         "",
					CommitStatus: false,
				},
				expected: "2_down",
			},
			{
				rawMigration: sqlmigrator.MigrationEntry{
					Number:       3,
					Name:         "create_posts_table",
					CommitStatus: true,
				},
				expected: "3_create_posts_table_up",
			},
			{
				rawMigration: sqlmigrator.MigrationEntry{
					Number:       4,
					Name:         "",
					CommitStatus: true,
				},
				expected: "4_up",
			},
		}

		for _, c := range tc {
			assert.Equal(t, c.expected, c.rawMigration.String())
		}
	})
}
