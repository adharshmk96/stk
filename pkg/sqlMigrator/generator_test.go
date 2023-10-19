package sqlmigrator_test

import (
	"fmt"
	"os"
	"path"
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

func getFileContent(t *testing.T, filePath string) string {
	t.Helper()
	file, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	return string(file)
}

func TestGenerate(t *testing.T) {
	t.Run("generator generates correct number of migrations", func(t *testing.T) {
		numToGenerate := 3
		expectedNumFiles := (numToGenerate * 2) + 1

		tempFolder := t.TempDir()

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(tempFolder, db, logfile, false)
		defer remoteFolder(tempFolder)

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, false)

		generatedFiles, err := generator.Generate(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedNumFiles, getNumberOfFilesInFolder(t, tempFolder))
		assert.Equal(t, expectedNumFiles-1, len(generatedFiles))

	})

	t.Run("generator fills file with content on fill flag", func(t *testing.T) {
		numToGenerate := 3

		tempFolder := t.TempDir()

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(tempFolder, db, logfile, false)
		defer remoteFolder(tempFolder)

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, true)

		generatedFiles, err := generator.Generate(ctx)

		assert.NoError(t, err)

		for _, file := range generatedFiles {
			content := getFileContent(t, file)
			assert.NotEmpty(t, content)
		}

	})

	t.Run("generator doesn't generate files on dry run", func(t *testing.T) {
		numToGenerate := 3

		tempFolder := t.TempDir()

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(tempFolder, db, logfile, true)
		defer remoteFolder(tempFolder)

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, true)

		generatedFiles, err := generator.Generate(ctx)

		assert.NoError(t, err)
		assert.Empty(t, generatedFiles)
		assert.Equal(t, 1, getNumberOfFilesInFolder(t, tempFolder))
	})

	t.Run("generator writes to log file", func(t *testing.T) {
		numToGenerate := 3
		tempFolder := t.TempDir()
		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewMigratorContext(tempFolder, db, logfile, false)

		logFilePath := path.Join(ctx.WorkDir, ctx.LogFile)

		defer remoteFolder(tempFolder)

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, true)
		_, err := generator.Generate(ctx)
		assert.NoError(t, err)

		logContent := getFileContent(t, logFilePath)
		assert.NotEmpty(t, logContent)

		expectedLogContent := func() string {
			content := ""
			for i := 1; i <= numToGenerate; i++ {
				content += fmt.Sprintf("%d_user_table\n", i)
			}
			return content
		}()
		assert.Equal(t, expectedLogContent, logContent)

		generator = sqlmigrator.NewGenerator("auth_table", numToGenerate+1, true)
		_, err = generator.Generate(ctx)
		assert.NoError(t, err)

		logContent = getFileContent(t, logFilePath)
		assert.NotEmpty(t, logContent)

		expectedLogContent = func() string {
			content := ""
			for i := 1; i <= numToGenerate; i++ {
				content += fmt.Sprintf("%d_user_table\n", i)
			}
			for i := numToGenerate + 1; i <= (2*numToGenerate)+1; i++ {
				content += fmt.Sprintf("%d_auth_table\n", i)
			}
			return content
		}()

		assert.Equal(t, expectedLogContent, logContent)
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
