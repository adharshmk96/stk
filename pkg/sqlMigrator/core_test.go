package sqlmigrator_test

import (
	"os"
	"path"
	"strings"
	"testing"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

var LOG_FILE_CONTENT = `1_create_users_table_up
2_create_posts_table_up
3_create_comments_table_up
4_create_likes_table_down
5_create_followers_table_down
6_create_messages_table_down
`

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
			assert.Equal(t, c.commitStatus, rawMigration.Committed)
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
					Number:    1,
					Name:      "create_users_table",
					Committed: true,
				},
				expected: "1_create_users_table_up",
			},
			{
				rawMigration: sqlmigrator.MigrationEntry{
					Number:    2,
					Name:      "",
					Committed: false,
				},
				expected: "2_down",
			},
			{
				rawMigration: sqlmigrator.MigrationEntry{
					Number:    3,
					Name:      "create_posts_table",
					Committed: true,
				},
				expected: "3_create_posts_table_up",
			},
			{
				rawMigration: sqlmigrator.MigrationEntry{
					Number:    4,
					Name:      "",
					Committed: true,
				},
				expected: "4_up",
			},
		}

		for _, c := range tc {
			assert.Equal(t, c.expected, c.rawMigration.EntryString())
		}
	})
}

func TestContextLoding(t *testing.T) {
	t.Run("loads migration entries from log file", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()

		ctx := sqlmigrator.NewMigratorContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)
		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)

		err := os.WriteFile(logFilePath, []byte(LOG_FILE_CONTENT), 0644)
		assert.NoError(t, err)

		err = ctx.LoadMigrationEntries()
		assert.NoError(t, err)

		assert.Equal(t, 6, len(ctx.Migrations))

		expected := func() []*sqlmigrator.MigrationEntry {
			migrationEntry := []*sqlmigrator.MigrationEntry{}

			lines := strings.Split(LOG_FILE_CONTENT, "\n")[0:6]
			for _, line := range lines {
				entry, err := sqlmigrator.ParseMigrationEntry(line)
				assert.NoError(t, err)
				migrationEntry = append(migrationEntry, entry)
			}

			return migrationEntry
		}()

		for i, migration := range ctx.Migrations {
			assert.Equal(t, expected[i].EntryString(), migration.EntryString())
			assert.Equal(t, expected[i].Number, migration.Number)
			assert.Equal(t, expected[i].Name, migration.Name)
			assert.Equal(t, expected[i].Committed, migration.Committed)
		}
	})

	t.Run("empty migration entries if log file is empty", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()

		ctx := sqlmigrator.NewMigratorContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)
		err := ctx.LoadMigrationEntries()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(ctx.Migrations))
	})

	t.Run("writes migration entries to logfile", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)
		defer removeDir()

		ctx := sqlmigrator.NewMigratorContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)
		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)

		err := ctx.LoadMigrationEntries()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(ctx.Migrations))

		migrationEntry := []*sqlmigrator.MigrationEntry{}

		lines := strings.Split(LOG_FILE_CONTENT, "\n")[0:6]
		for _, line := range lines {
			entry, err := sqlmigrator.ParseMigrationEntry(line)
			assert.NoError(t, err)
			migrationEntry = append(migrationEntry, entry)
		}

		ctx.Migrations = migrationEntry

		err = ctx.WriteMigrationEntries()
		assert.NoError(t, err)

		assert.FileExists(t, logFilePath)
		assert.Equal(t, LOG_FILE_CONTENT, testutils.GetFileContent(t, logFilePath))
	})
}
