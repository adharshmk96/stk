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

func getNumberOfFilesInFolder(t *testing.T, folder string) int {
	t.Helper()
	files, err := os.ReadDir(folder)
	assert.NoError(t, err)
	return len(files)
}

func TestGenerate(t *testing.T) {
	t.Run("generator generates correct number of migrations", func(t *testing.T) {
		numToGenerate := 3
		expectedNumFiles := (numToGenerate * 2) + 1

		tempDir, removeDir := testutils.CreateTempDirectory(t)

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(tempDir, db, logfile, false)
		defer removeDir()

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, false)

		generatedFiles, err := generator.Generate(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedNumFiles, getNumberOfFilesInFolder(t, tempDir))
		assert.Equal(t, expectedNumFiles-1, len(generatedFiles))

	})

	t.Run("generator fills file with content on fill flag", func(t *testing.T) {
		numToGenerate := 3

		tempDir, removeDir := testutils.CreateTempDirectory(t)

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(tempDir, db, logfile, false)
		defer removeDir()

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, true)

		generatedFiles, err := generator.Generate(ctx)

		assert.NoError(t, err)

		for _, file := range generatedFiles {
			content := testutils.GetFileContent(t, file)
			assert.NotEmpty(t, content)
		}

	})

	t.Run("generator doesn't generate files on dry run", func(t *testing.T) {
		numToGenerate := 3

		tempDir, removeDir := testutils.CreateTempDirectory(t)

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(tempDir, db, logfile, true)
		defer removeDir()

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, true)

		generatedFiles, err := generator.Generate(ctx)

		assert.NoError(t, err)
		assert.Empty(t, generatedFiles)
		assert.Equal(t, 1, getNumberOfFilesInFolder(t, tempDir))
	})

	t.Run("generator writes to log file", func(t *testing.T) {
		numToGenerate := 3

		tempDir, removeDir := testutils.CreateTempDirectory(t)

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(tempDir, db, logfile, false)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)

		defer removeDir()

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, true)
		_, err := generator.Generate(ctx)
		assert.NoError(t, err)

		logContent := testutils.GetFileContent(t, logFilePath)
		assert.NotEmpty(t, logContent)

		expectedLogContent := func() string {
			content := ""
			for i := 1; i <= numToGenerate; i++ {
				content += fmt.Sprintf("%d_user_table_down\n", i)
			}
			return content
		}()
		assert.Equal(t, expectedLogContent, logContent)

		generator = sqlmigrator.NewGenerator("auth_table", numToGenerate+1, true)
		_, err = generator.Generate(ctx)
		assert.NoError(t, err)

		logContent = testutils.GetFileContent(t, logFilePath)
		assert.NotEmpty(t, logContent)

		expectedLogContent = func() string {
			content := ""
			for i := 1; i <= numToGenerate; i++ {
				content += fmt.Sprintf("%d_user_table_down\n", i)
			}
			for i := numToGenerate + 1; i <= (2*numToGenerate)+1; i++ {
				content += fmt.Sprintf("%d_auth_table_down\n", i)
			}
			return content
		}()

		assert.Equal(t, expectedLogContent, logContent)
	})
}

func TestGenerateNextMigrations(t *testing.T) {
	t.Run("generates next migrations", func(t *testing.T) {
		lastMigration := sqlmigrator.MigrationEntry{
			Number: 1,
			Name:   "create_users_table",
		}

		migrations := []sqlmigrator.MigrationEntry{
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
