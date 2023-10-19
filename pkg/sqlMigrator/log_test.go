package sqlmigrator_test

import (
	"os"
	"path"
	"testing"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/stretchr/testify/assert"
)

func InitializeMigrationsFolder(t *testing.T) {
	t.Run("creates a migrations folder", func(t *testing.T) {
		ctx := sqlmigrator.NewMigratorContext(t.TempDir(), sqlmigrator.SQLiteDB, "migrator.log", false)
		err := sqlmigrator.InitializeMigrationsFolder(ctx)
		assert.NoError(t, err)

		_, err = os.Stat(path.Join(ctx.WorkDir, ctx.LogFile))
		assert.NoError(t, err)
	})

	t.Run("does not create a more log file if it already exists", func(t *testing.T) {
		ctx := sqlmigrator.NewMigratorContext(t.TempDir(), sqlmigrator.SQLiteDB, "migrator.log", false)
		logPath := path.Join(ctx.WorkDir, ctx.LogFile)

		err := sqlmigrator.InitializeMigrationsFolder(ctx)
		assert.NoError(t, err)

		numFiles := getNumberOfFilesInFolder(t, logPath)
		assert.Equal(t, 1, numFiles)

		err = sqlmigrator.InitializeMigrationsFolder(ctx)
		assert.NoError(t, err)

		numFiles = getNumberOfFilesInFolder(t, logPath)
		assert.Equal(t, 1, numFiles)

	})
}
