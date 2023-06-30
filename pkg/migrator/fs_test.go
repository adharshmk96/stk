package migrator_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/adharshmk96/stk/pkg/migrator"
)

var upFileNames = []string{
	"1_up",
	"2_up",
	"000004_up",
	"000005_up",
	"1000001_create_table_up",
	"1000002_create_table2_up",
}

var downFileNames = []string{
	"1_down",
	"2_down",
	"000004_down",
	"000005_down",
	"1000001_create_table_down",
	"1000002_create_table2_down",
}

var noiseFiles = []string{
	"1.txt",
	"2.txt",
	"000004.txt",
	"whatever.txt",
}

var testDir = "test_dir"

func setupDir() {
	os.MkdirAll(testDir, 0700)

	// create all upFileNames
	for _, name := range upFileNames {
		f, err := os.Create(filepath.Join(testDir, name+".sql"))
		if err != nil {
			fmt.Println(err)
		}
		f.Close()
	}

	// create all downFileNames
	for _, name := range downFileNames {
		f, err := os.Create(filepath.Join(testDir, name+".sql"))
		if err != nil {
			fmt.Println(err)
		}
		f.Close()
	}

	// create all noiseFiles
	for _, name := range noiseFiles {
		f, err := os.Create(filepath.Join(testDir, name+".sql"))
		if err != nil {
			fmt.Println(err)
		}
		f.Close()
	}

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
	setupDir()
	defer teardownDir()
	t.Run("get up filenames without extension", func(t *testing.T) {
		filenames, err := migrator.GetMigrationFileGroup("test_dir", migrator.MigrationUp)
		if err != nil {
			t.Error(err)
		}
		for _, name := range upFileNames {
			if !contains(filenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}

		for _, name := range downFileNames {
			if contains(filenames, name) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}

		for _, name := range noiseFiles {
			if contains(filenames, name) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}
	})

	t.Run("get down filenames without extension", func(t *testing.T) {
		filenames, err := migrator.GetMigrationFileGroup("test_dir", migrator.MigrationDown)
		if err != nil {
			t.Error(err)
		}
		for _, name := range downFileNames {
			if !contains(filenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}

		for _, name := range upFileNames {
			if contains(filenames, name) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}

		for _, name := range noiseFiles {
			if contains(filenames, name) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}
	})

}
