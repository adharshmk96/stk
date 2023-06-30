package migrator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/adharshmk96/stk/pkg/migrator/testutils"
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

func setupFSDir() {
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

func teardownFSDir() {
	os.RemoveAll("test_dir")
}

func TestGetFilenamesWithoutExtension(t *testing.T) {
	setupFSDir()
	defer teardownFSDir()
	t.Run("get up filenames without extension", func(t *testing.T) {
		filenames, err := getMigrationFileGroup("test_dir", MigrationUp)
		if err != nil {
			t.Error(err)
		}
		for _, name := range upFileNames {
			if !testutils.Contains(filenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}

		for _, name := range downFileNames {
			if testutils.Contains(filenames, name) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}

		for _, name := range noiseFiles {
			if testutils.Contains(filenames, name) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}
	})

	t.Run("get down filenames without extension", func(t *testing.T) {
		filenames, err := getMigrationFileGroup("test_dir", MigrationDown)
		if err != nil {
			t.Error(err)
		}
		for _, name := range downFileNames {
			if !testutils.Contains(filenames, name) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}

		for _, name := range upFileNames {
			if testutils.Contains(filenames, name) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}

		for _, name := range noiseFiles {
			if testutils.Contains(filenames, name) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}
	})

}
