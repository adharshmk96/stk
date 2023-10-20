package sqlmigrator_test

import (
	"os"
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
		ctx := sqlmigrator.NewContext(tempDir, db, logfile, false)
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
		ctx := sqlmigrator.NewContext(tempDir, db, logfile, false)
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
		ctx := sqlmigrator.NewContext(tempDir, db, logfile, true)
		defer removeDir()

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, true)

		generatedFiles, err := generator.Generate(ctx)

		assert.NoError(t, err)
		assert.Empty(t, generatedFiles)
		assert.Equal(t, 1, getNumberOfFilesInFolder(t, tempDir))
	})

	t.Run("generator updates ctx migrations", func(t *testing.T) {
		numToGenerate := 3

		tempDir, removeDir := testutils.CreateTempDirectory(t)

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(tempDir, db, logfile, false)

		defer removeDir()

		existingMigrations := len(ctx.Migrations)

		generator := sqlmigrator.NewGenerator("user_table", numToGenerate, true)
		_, err := generator.Generate(ctx)
		assert.NoError(t, err)

		assert.Equal(t, existingMigrations+numToGenerate, len(ctx.Migrations))

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

func TestClean(t *testing.T) {
	t.Run("clean removes files and updates ctx", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(tempDir, db, logfile, false)
		defer removeDir()

		generator := sqlmigrator.NewGenerator("user_table", 3, true)
		generatedFiles, err := generator.Generate(ctx)
		assert.NoError(t, err)

		for _, migration := range ctx.Migrations {
			migration.Committed = true
		}

		expectedNumFiles := 1 + len(generatedFiles)
		exptctedNumMigrations := (expectedNumFiles - 1) / 2

		assert.Equal(t, expectedNumFiles, getNumberOfFilesInFolder(t, tempDir))
		assert.Equal(t, exptctedNumMigrations, len(ctx.Migrations))

		generator = sqlmigrator.NewGenerator("groups_table", 4, true)
		uncommitedFiles, err := generator.Generate(ctx)
		assert.NoError(t, err)

		expectedNumFiles = 1 + len(generatedFiles) + len(uncommitedFiles)
		exptctedNumMigrations = (expectedNumFiles - 1) / 2

		assert.Equal(t, expectedNumFiles, getNumberOfFilesInFolder(t, tempDir))
		assert.Equal(t, exptctedNumMigrations, len(ctx.Migrations))

		for _, migration := range ctx.Migrations {
			assert.FileExists(t, migration.UpFilePath)
			assert.FileExists(t, migration.DownFilePath)
		}

		err = ctx.WriteMigrationEntries()
		assert.NoError(t, err)

		removedFiles, err := generator.Clean(ctx)
		assert.NoError(t, err)

		assert.Equal(t, len(uncommitedFiles), len(removedFiles))

		expectedNumFiles = 1 + len(generatedFiles)
		exptctedNumMigrations = (expectedNumFiles - 1) / 2

		assert.Equal(t, expectedNumFiles, getNumberOfFilesInFolder(t, tempDir))
		assert.Equal(t, exptctedNumMigrations, len(ctx.Migrations))

		for _, migration := range ctx.Migrations {
			assert.FileExists(t, migration.UpFilePath)
			assert.FileExists(t, migration.DownFilePath)
		}

	})

	t.Run("clean doesn't remove files on dry run", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(tempDir, db, logfile, false)
		defer removeDir()

		generator := sqlmigrator.NewGenerator("user_table", 3, true)
		generatedFiles, err := generator.Generate(ctx)
		assert.NoError(t, err)

		for _, migration := range ctx.Migrations {
			assert.FileExists(t, migration.UpFilePath)
			assert.FileExists(t, migration.DownFilePath)
			migration.Committed = true
		}

		generator = sqlmigrator.NewGenerator("groups_table", 4, true)
		uncommitedFiles, err := generator.Generate(ctx)
		assert.NoError(t, err)

		ctx.DryRun = true

		removedFiles, err := generator.Clean(ctx)
		assert.NoError(t, err)

		expectedFiles := len(uncommitedFiles) + len(generatedFiles) + 1
		expectedMigrations := (len(uncommitedFiles) + len(generatedFiles)) / 2

		assert.Empty(t, removedFiles)
		assert.Equal(t, expectedFiles, getNumberOfFilesInFolder(t, tempDir))
		assert.Equal(t, expectedMigrations, len(ctx.Migrations))
	})

	t.Run("clean works in empty directory", func(t *testing.T) {
		tempDir, removeDir := testutils.CreateTempDirectory(t)

		_, db, logfile := sqlmigrator.DefaultContextConfig()
		ctx := sqlmigrator.NewContext(tempDir, db, logfile, false)
		defer removeDir()

		generator := sqlmigrator.NewGenerator("user_table", 3, true)
		removedFiles, err := generator.Clean(ctx)
		assert.NoError(t, err)

		assert.Empty(t, removedFiles)
		assert.Equal(t, 1, getNumberOfFilesInFolder(t, tempDir))
		assert.Equal(t, 0, len(ctx.Migrations))
	})
}
