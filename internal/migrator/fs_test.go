package migrator_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/adharshmk96/stk/internal/migrator"
)

func setupDir() []string {
	testDir := "test_dir"
	os.MkdirAll(testDir, 0700)

	fileNames := []string{}
	for i := 0; i < 10; i++ {
		fileName := "test_file" + fmt.Sprintf("%d", i)
		fileNames = append(fileNames, fileName)
		path := filepath.Join(testDir, fileName+".sql")
		os.Create(path)
	}
	return fileNames
}

func teardownDir() {
	os.RemoveAll("test_dir")
}

func contains(s []string, e string) bool {
	for _, name := range s {
		if name == e {
			return true
		}
	}
	return false
}

func TestGetFilenamesWithoutExtension(t *testing.T) {
	existingNames := setupDir()
	defer teardownDir()
	t.Run("get filenames without extension", func(t *testing.T) {
		filenames, err := migrator.GetFilenamesWithoutExtension("test_dir")
		if err != nil {
			t.Error(err)
		}
		for _, name := range existingNames {
			if !contains(filenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}
	})
}
