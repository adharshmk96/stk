package migrator_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/internal/migrator"
	"github.com/stretchr/testify/assert"
)

func teardownGenDir() {
	os.RemoveAll("test_dir")
	os.RemoveAll("test_dir2")
}

func TestGenerator(t *testing.T) {

	var testDir = "test_dir"
	defer teardownGenDir()

	t.Run("generator generates files in empty directory", func(t *testing.T) {
		expectedUpFileNames := []string{
			"000001_test_up",
			"000002_test_up",
			"000003_test_up",
			"000004_test_up",
		}
		expectedDownFileNames := []string{
			"000001_test_down",
			"000002_test_down",
			"000003_test_down",
			"000004_test_down",
		}
		config := migrator.GeneratorConfig{
			RootDirectory: testDir,
			Database:      "sqlite",
			Name:          "test",
			NumToGenerate: 4,
			DryRun:        false,
		}

		err := migrator.Generate(config)
		if err != nil {
			t.Error(err)
		}
		upFilenames, err := migrator.GetMigrationFileGroup("test_dir/sqlite", migrator.MigrationUp)
		assert.NoError(t, err)
		downFilenames, err := migrator.GetMigrationFileGroup("test_dir/sqlite", migrator.MigrationDown)
		assert.NoError(t, err)

		assert.Equal(t, 4, len(upFilenames))
		assert.Equal(t, 4, len(downFilenames))

		for _, name := range expectedUpFileNames {
			if !contains(upFilenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}
		for _, name := range expectedDownFileNames {
			if !contains(downFilenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}
	})

	t.Run("generator generates files in non-empty directory", func(t *testing.T) {
		expectedUpFileNames := []string{
			"000004_test_up",
			"000005_test_up",
			"000006_test_up",
			"000007_test_up",
		}
		expectedDownFileNames := []string{
			"000004_test_down",
			"000005_test_down",
			"000006_test_down",
			"000007_test_down",
		}
		config := migrator.GeneratorConfig{
			RootDirectory: "test_dir2",
			Database:      "sqlite",
			Name:          "test",
			NumToGenerate: 4,
			DryRun:        false,
		}

		_ = migrator.Generate(config)
		err := migrator.Generate(config)
		if err != nil {
			t.Error(err)
		}
		upFilenames, err := migrator.GetMigrationFileGroup("test_dir2/sqlite", migrator.MigrationUp)
		assert.NoError(t, err)
		downFilenames, err := migrator.GetMigrationFileGroup("test_dir2/sqlite", migrator.MigrationDown)
		assert.NoError(t, err)

		assert.Equal(t, 8, len(upFilenames))
		assert.Equal(t, 8, len(downFilenames))

		for _, name := range expectedUpFileNames {
			if !contains(upFilenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}
		for _, name := range expectedDownFileNames {
			if !contains(downFilenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}
	})
}
