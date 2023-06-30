package migrator

import (
	"fmt"
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

		start := 0
		numToGenerate := 4
		expectedNumberOfFiles := start + numToGenerate

		var expectedUpFileNames []string
		var expectedDownFileNames []string

		for i := start + 1; i <= start+numToGenerate; i++ {
			up := fmt.Sprintf("%06d_test_up", i)
			down := fmt.Sprintf("%06d_test_down", i)
			expectedUpFileNames = append(expectedUpFileNames, up)
			expectedDownFileNames = append(expectedDownFileNames, down)
		}

		config := GeneratorConfig{
			RootDirectory: testDir,
			Database:      "sqlite",
			Name:          "test",
			NumToGenerate: numToGenerate,
			DryRun:        false,
		}

		err := Generate(config)
		if err != nil {
			t.Error(err)
		}
		upFilenames, err := getMigrationFileNamesByGroup("test_dir/sqlite", MigrationUp)
		assert.NoError(t, err)
		downFilenames, err := getMigrationFileNamesByGroup("test_dir/sqlite", MigrationDown)
		assert.NoError(t, err)

		assert.Equal(t, expectedNumberOfFiles, len(upFilenames))
		assert.Equal(t, expectedNumberOfFiles, len(downFilenames))

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
		start := 4
		numToGenerate := 4
		expectedNumberOfFiles := start + numToGenerate

		var expectedUpFileNames []string
		var expectedDownFileNames []string

		for i := start + 1; i <= start+numToGenerate; i++ {
			up := fmt.Sprintf("%06d_test_up", i)
			down := fmt.Sprintf("%06d_test_down", i)
			expectedUpFileNames = append(expectedUpFileNames, up)
			expectedDownFileNames = append(expectedDownFileNames, down)
		}
		config := GeneratorConfig{
			RootDirectory: "test_dir2",
			Database:      "sqlite",
			Name:          "test",
			NumToGenerate: numToGenerate,
			DryRun:        false,
		}

		_ = Generate(config)
		err := Generate(config)
		if err != nil {
			t.Error(err)
		}
		upFilenames, err := getMigrationFileNamesByGroup("test_dir2/sqlite", MigrationUp)
		assert.NoError(t, err)
		downFilenames, err := getMigrationFileNamesByGroup("test_dir2/sqlite", MigrationDown)
		assert.NoError(t, err)

		assert.Equal(t, expectedNumberOfFiles, len(upFilenames))
		assert.Equal(t, expectedNumberOfFiles, len(downFilenames))

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
