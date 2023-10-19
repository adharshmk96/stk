package sqlmigrator_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func TestInitializeMigrationsFolder(t *testing.T) {
	t.Run("creates a migrations folder", func(t *testing.T) {
		ctx := sqlmigrator.NewMigratorContext(t.TempDir(), sqlmigrator.SQLiteDB, "migrator.log", false)
		err := sqlmigrator.InitializeMigrationsFolder(ctx)
		assert.NoError(t, err)

		_, err = os.Stat(path.Join(ctx.WorkDir, ctx.LogFile))
		assert.NoError(t, err)
	})

	t.Run("does not create a more log file if it already exists", func(t *testing.T) {

		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewMigratorContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)

		err := sqlmigrator.InitializeMigrationsFolder(ctx)
		assert.NoError(t, err)

		numFiles := getNumberOfFilesInFolder(t, tempDir)
		assert.Equal(t, 1, numFiles)

		err = sqlmigrator.InitializeMigrationsFolder(ctx)
		assert.NoError(t, err)

		numFiles = getNumberOfFilesInFolder(t, tempDir)
		assert.Equal(t, 1, numFiles)

	})
}

func TestLoadUncommitedMigrationFromLog(t *testing.T) {
	t.Run("returns an empty migration entry if the log file is empty", func(t *testing.T) {
		ctx := sqlmigrator.NewMigratorContext(t.TempDir(), sqlmigrator.SQLiteDB, "migrator.log", false)
		migration, err := sqlmigrator.LoadUncommitedMigrations(ctx)
		assert.NoError(t, err)
		assert.Empty(t, migration)
		assert.Equal(t, 0, len(migration))
	})

	t.Run("returns uncommitted migration entries from the log file", func(t *testing.T) {
		logFile_content := func() string {
			fileContent := ""
			for i := 1; i <= 3; i++ {
				fileContent += fmt.Sprintf("%d_create_users_table_up\n", i)
			}
			for i := 4; i <= 6; i++ {
				fileContent += fmt.Sprintf("%d_create_users_table_down\n", i)
			}
			return fileContent
		}()

		ctx := sqlmigrator.NewMigratorContext(t.TempDir(), sqlmigrator.SQLiteDB, "migrator.log", false)
		logPath := path.Join(ctx.WorkDir, ctx.LogFile)
		err := os.WriteFile(logPath, []byte(logFile_content), 0644)
		assert.NoError(t, err)

		migrations, err := sqlmigrator.LoadUncommitedMigrations(ctx)
		assert.NoError(t, err)

		expected := func() []*sqlmigrator.MigrationEntry {
			migrationEntry := []*sqlmigrator.MigrationEntry{}
			for i := 4; i <= 6; i++ {
				migrationEntry = append(migrationEntry, &sqlmigrator.MigrationEntry{
					Number:    i,
					Name:      "create_users_table",
					Committed: false,
				})
			}

			return migrationEntry
		}()

		for i, migration := range migrations {
			assert.Equal(t, expected[i].Name, migration.Name)
			assert.Equal(t, expected[i].Number, migration.Number)
			assert.Equal(t, expected[i].Committed, migration.Committed)
		}

	})
}
