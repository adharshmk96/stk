package migrator

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/pkg/migrator/testutils"
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
		config := GeneratorConfig{
			RootDirectory: testDir,
			Database:      "sqlite",
			Name:          "test",
			NumToGenerate: 4,
			DryRun:        false,
		}

		err := Generate(config)
		if err != nil {
			t.Error(err)
		}
		upFilenames, err := getMigrationFileGroup("test_dir/sqlite", MigrationUp)
		assert.NoError(t, err)
		downFilenames, err := getMigrationFileGroup("test_dir/sqlite", MigrationDown)
		assert.NoError(t, err)

		assert.Equal(t, 4, len(upFilenames))
		assert.Equal(t, 4, len(downFilenames))

		for _, name := range expectedUpFileNames {
			if !testutils.Contains(upFilenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}
		for _, name := range expectedDownFileNames {
			if !testutils.Contains(downFilenames, name) {
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
		config := GeneratorConfig{
			RootDirectory: "test_dir2",
			Database:      "sqlite",
			Name:          "test",
			NumToGenerate: 4,
			DryRun:        false,
		}

		_ = Generate(config)
		err := Generate(config)
		if err != nil {
			t.Error(err)
		}
		upFilenames, err := getMigrationFileGroup("test_dir2/sqlite", MigrationUp)
		assert.NoError(t, err)
		downFilenames, err := getMigrationFileGroup("test_dir2/sqlite", MigrationDown)
		assert.NoError(t, err)

		assert.Equal(t, 8, len(upFilenames))
		assert.Equal(t, 8, len(downFilenames))

		for _, name := range expectedUpFileNames {
			if !testutils.Contains(upFilenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}
		for _, name := range expectedDownFileNames {
			if !testutils.Contains(downFilenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}
	})
}
