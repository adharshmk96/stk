package sqlmigrator_test

import (
	"os"
	"path"
	"testing"

	"github.com/adharshmk96/stk/mocks"
	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func checkUnappliedMigrations(t *testing.T, ctx *sqlmigrator.Context, expected int) {
	unappliedMigrations, err := sqlmigrator.LoadUncommitedMigrations(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, len(unappliedMigrations))
}

func TestMigrateUp(t *testing.T) {

	var LOG_FILE_CONTENT = `1_create_users_table_up
2_create_posts_table_up
3_create_comments_table_up
4_create_likes_table_down
5_create_followers_table_down
6_create_messages_table_down
`

	t.Run("migrate up default all unapplied migrations", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)
		err := os.WriteFile(logFilePath, []byte(LOG_FILE_CONTENT), 0644)
		assert.NoError(t, err)
		err = ctx.LoadMigrationEntries()
		assert.NoError(t, err)
		checkUnappliedMigrations(t, ctx, 3)

		dbMock := mocks.NewDBRepo(t)
		dbMock.On("Exec", mock.AnythingOfType("string")).Return(nil)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateUp(ctx, 0)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(appliedMigrations))

		expected := []string{
			"4_create_likes_table_up",
			"5_create_followers_table_up",
			"6_create_messages_table_up",
		}

		for i, migration := range appliedMigrations {
			assert.True(t, migration.Committed)
			assert.Equal(t, expected[i], migration.EntryString())
		}

		checkUnappliedMigrations(t, ctx, 0)
	})

	t.Run("migrate up given number of unapplied migration", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)
		err := os.WriteFile(logFilePath, []byte(LOG_FILE_CONTENT), 0644)
		assert.NoError(t, err)
		err = ctx.LoadMigrationEntries()
		assert.NoError(t, err)
		checkUnappliedMigrations(t, ctx, 3)

		dbMock := mocks.NewDBRepo(t)
		dbMock.On("Exec", mock.AnythingOfType("string")).Return(nil)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateUp(ctx, 1)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(appliedMigrations))

		expected := []string{
			"4_create_likes_table_up",
		}

		for i, migration := range appliedMigrations {
			assert.True(t, migration.Committed)
			assert.Equal(t, expected[i], migration.EntryString())
		}

		// Check unapplied migrations
		checkUnappliedMigrations(t, ctx, 2)
	})

	t.Run("migrate up won't update commit for dry run", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", true)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)
		err := os.WriteFile(logFilePath, []byte(LOG_FILE_CONTENT), 0644)
		assert.NoError(t, err)
		err = ctx.LoadMigrationEntries()
		assert.NoError(t, err)
		checkUnappliedMigrations(t, ctx, 3)

		dbMock := mocks.NewDBRepo(t)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateUp(ctx, 0)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(appliedMigrations))

		checkUnappliedMigrations(t, ctx, 3)

	})

	t.Run("migrate up runs fine with outof bound values", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)
		err := os.WriteFile(logFilePath, []byte(LOG_FILE_CONTENT), 0644)
		assert.NoError(t, err)
		err = ctx.LoadMigrationEntries()
		assert.NoError(t, err)
		checkUnappliedMigrations(t, ctx, 3)

		dbMock := mocks.NewDBRepo(t)
		dbMock.On("Exec", mock.AnythingOfType("string")).Return(nil)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateUp(ctx, 0)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(appliedMigrations))

		checkUnappliedMigrations(t, ctx, 0)
	})

	t.Run("migrate up runs fine with no unapplied migrations", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)

		checkUnappliedMigrations(t, ctx, 0)

		dbMock := mocks.NewDBRepo(t)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateUp(ctx, 0)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(appliedMigrations))

		checkUnappliedMigrations(t, ctx, 0)
	})
}

func TestMigrateDown(t *testing.T) {

	var LOG_FILE_CONTENT = `1_create_users_table_up
2_create_posts_table_up
3_create_comments_table_up
4_create_likes_table_down
5_create_followers_table_down
6_create_messages_table_down
`

	t.Run("migrate down default all unapplied migrations", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)
		err := os.WriteFile(logFilePath, []byte(LOG_FILE_CONTENT), 0644)
		assert.NoError(t, err)

		err = ctx.LoadMigrationEntries()
		assert.NoError(t, err)

		checkUnappliedMigrations(t, ctx, 3)

		dbMock := mocks.NewDBRepo(t)
		dbMock.On("Exec", mock.AnythingOfType("string")).Return(nil)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateDown(ctx, 0)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(appliedMigrations))

		expected := []string{
			"3_create_comments_table_down",
			"2_create_posts_table_down",
			"1_create_users_table_down",
		}

		for i, migration := range appliedMigrations {
			assert.False(t, migration.Committed)
			assert.Equal(t, expected[i], migration.EntryString())
		}

		checkUnappliedMigrations(t, ctx, 6)
	})

	t.Run("migrate down given number of unapplied migration", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)
		err := os.WriteFile(logFilePath, []byte(LOG_FILE_CONTENT), 0644)
		assert.NoError(t, err)

		err = ctx.LoadMigrationEntries()
		assert.NoError(t, err)

		checkUnappliedMigrations(t, ctx, 3)

		dbMock := mocks.NewDBRepo(t)
		dbMock.On("Exec", mock.AnythingOfType("string")).Return(nil)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateDown(ctx, 1)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(appliedMigrations))

		expected := []string{
			"3_create_comments_table_down",
		}

		for i, migration := range appliedMigrations {
			assert.False(t, migration.Committed)
			assert.Equal(t, expected[i], migration.EntryString())
		}

		// Check unapplied migrations
		checkUnappliedMigrations(t, ctx, 4)
	})

	t.Run("migrate down won't update commit for dry run", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", true)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)
		err := os.WriteFile(logFilePath, []byte(LOG_FILE_CONTENT), 0644)
		assert.NoError(t, err)

		err = ctx.LoadMigrationEntries()
		assert.NoError(t, err)

		checkUnappliedMigrations(t, ctx, 3)

		dbMock := mocks.NewDBRepo(t)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateDown(ctx, 0)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(appliedMigrations))

		checkUnappliedMigrations(t, ctx, 3)
	})

	t.Run("migrate down runs fine with outof bound values", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)
		err := os.WriteFile(logFilePath, []byte(LOG_FILE_CONTENT), 0644)
		assert.NoError(t, err)

		err = ctx.LoadMigrationEntries()
		assert.NoError(t, err)

		checkUnappliedMigrations(t, ctx, 3)

		dbMock := mocks.NewDBRepo(t)
		dbMock.On("Exec", mock.AnythingOfType("string")).Return(nil)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateDown(ctx, 0)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(appliedMigrations))

		checkUnappliedMigrations(t, ctx, 6)
	})

	t.Run("migrate down runs fine with no unapplied migrations", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		defer removeDir()

		ctx := sqlmigrator.NewContext(tempDir, sqlmigrator.SQLiteDB, "migrator.log", false)

		checkUnappliedMigrations(t, ctx, 0)

		dbMock := mocks.NewDBRepo(t)

		migrator := sqlmigrator.NewMigrator(dbMock)
		appliedMigrations, err := migrator.MigrateDown(ctx, 0)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(appliedMigrations))

		checkUnappliedMigrations(t, ctx, 0)
	})
}
