package sqlmigrator_test

import (
	"os"
	"testing"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	"github.com/stretchr/testify/assert"
)

func remoteFolder(folder string) {
	os.RemoveAll(folder)
}

func getNumberOfFilesInFolder(t *testing.T, folder string) int {
	t.Helper()
	files, err := os.ReadDir(folder)
	assert.NoError(t, err)
	return len(files)
}

func TestGenerate(t *testing.T) {
	t.Run("generator generates correct number of migrations", func(t *testing.T) {
		numToGenerate := 3
		expectedFiles := (numToGenerate * 2) + 1

		tempFolder := t.TempDir()

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(tempFolder, db, logfile, false)
		defer remoteFolder(tempFolder)

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, false, false)

		err := generator.Generate(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedFiles, getNumberOfFilesInFolder(t, tempFolder))

	})
}

func TestGenerateNextMigrations(t *testing.T) {
	t.Run("generates next migrations", func(t *testing.T) {
		lastMigration := sqlmigrator.Migration{
			Number: 1,
			Name:   "create_users_table",
		}

		migrations := []sqlmigrator.Migration{
			{
				Number: 2,
			},
			{
				Number: 3,
			},
		}

		nextMigrations := sqlmigrator.GenerateNextMigrations(lastMigration.Number, "", len(migrations))

		for i, migration := range nextMigrations {
			assert.Equal(t, migrations[i].Number, migration.Number)
		}

	})
}
