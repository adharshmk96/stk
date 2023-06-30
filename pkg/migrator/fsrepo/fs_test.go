package fsrepo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/adharshmk96/stk/pkg/migrator/testutils"
	"github.com/stretchr/testify/assert"
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
var emptyDir = "empty_dir"
var extention = migrator.GetExtention(migrator.SQLiteDB)

func setupFSDir() {
	os.MkdirAll(testDir, 0700)
	os.MkdirAll(emptyDir, 0700)

	// create all upFileNames
	for _, name := range upFileNames {
		f, err := os.Create(filepath.Join(testDir, name+"."+extention))
		if err != nil {
			fmt.Println(err)
		}
		f.Close()
	}

	// create all downFileNames
	for _, name := range downFileNames {
		f, err := os.Create(filepath.Join(testDir, name+"."+extention))
		if err != nil {
			fmt.Println(err)
		}
		f.Close()
	}

	// create all noiseFiles
	for _, name := range noiseFiles {
		f, err := os.Create(filepath.Join(testDir, name+"."+extention))
		if err != nil {
			fmt.Println(err)
		}
		f.Close()
	}

}

func teardownFSDir() {
	os.RemoveAll(testDir)
	os.RemoveAll(emptyDir)
}

func TestParseMigrationsFromFilePaths(t *testing.T) {
	t.Run("parses migrations from file paths", func(t *testing.T) {
		var filepaths []string
		expected := make([]*migrator.Migration, 0)

		for i := 1; i <= 10; i++ {
			name := fmt.Sprintf("create_users_table%d", i)

			up := fmt.Sprintf("%06d_%s_up.sql", i, name)
			down := fmt.Sprintf("%06d_%s_down.sql", i, name)
			filepaths = append(filepaths, up)
			filepaths = append(filepaths, down)

			expected = append(expected, &migrator.Migration{
				Number: i,
				Name:   name,
				Type:   migrator.MigrationUp,
				Path:   up,
			})
			expected = append(expected, &migrator.Migration{
				Number: i,
				Name:   name,
				Type:   migrator.MigrationDown,
				Path:   down,
			})
		}

		actual, err := parseMigrationsFromFilePaths(filepaths)

		assert.NoError(t, err)

		for i := range filepaths {
			assert.Equal(t, expected[i].Number, actual[i].Number)
			assert.Equal(t, expected[i].Name, actual[i].Name)
			assert.Equal(t, expected[i].Type, actual[i].Type)
			assert.Equal(t, expected[i].Path, actual[i].Path)
		}

	})

	t.Run("returns error if file path is invalid", func(t *testing.T) {
		filepaths := []string{
			"invalid",
		}

		_, err := parseMigrationsFromFilePaths(filepaths)

		assert.Error(t, err)
	})
}

func TestLoadMigrationsFromFile(t *testing.T) {

	t.Run("load migrations from empty directory", func(t *testing.T) {
		ext := migrator.GetExtention(migrator.SQLiteDB)
		fsRepo := NewFSRepo("test_dir", ext)
		migrations, err := fsRepo.LoadMigrationsFromFile(migrator.MigrationUp)
		assert.NoError(t, err)
		assert.Equal(t, len(migrations), 0)
	})
}

func TestGetMigrationFilePathsByGroup(t *testing.T) {
	setupFSDir()
	defer teardownFSDir()

	ext := migrator.GetExtention(migrator.SQLiteDB)
	fsRepo := NewFSRepo("test_dir", ext)

	t.Run("get up filenames with extension", func(t *testing.T) {
		filenames, err := fsRepo.GetMigrationFilePathsByType(migrator.MigrationUp)
		if err != nil {
			t.Error(err)
		}
		for _, name := range upFileNames {
			if !testutils.Contains(filenames, filepath.Join(testDir, name+"."+extention)) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}

		for _, name := range downFileNames {
			if testutils.Contains(filenames, filepath.Join(testDir, name+"."+extention)) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}

		for _, name := range noiseFiles {
			if testutils.Contains(filenames, filepath.Join(testDir, name+"."+extention)) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}
	})

	t.Run("get down filenames with extension", func(t *testing.T) {
		filenames, err := fsRepo.GetMigrationFilePathsByType(migrator.MigrationDown)
		if err != nil {
			t.Error(err)
		}
		for _, name := range downFileNames {
			if !testutils.Contains(filenames, filepath.Join(testDir, name+"."+extention)) {
				t.Errorf("expected %s to be in filenames", name)
			}
		}

		for _, name := range upFileNames {
			if testutils.Contains(filenames, filepath.Join(testDir, name+"."+extention)) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}

		for _, name := range noiseFiles {
			if testutils.Contains(filenames, filepath.Join(testDir, name+"."+extention)) {
				t.Errorf("expected %s to not be in filenames", name)
			}
		}
	})
}

func TestCreateMigrationFile(t *testing.T) {
	setupFSDir()
	defer teardownFSDir()

	ext := migrator.GetExtention(migrator.SQLiteDB)
	// subDir := migrator.SelectSubDirectory(migrator.SQLiteDB)
	fsRepo := NewFSRepo(testDir, ext)

	t.Run("create up migration file", func(t *testing.T) {
		migration := &migrator.Migration{
			Number: 100,
			Name:   "create_users_table",
			Type:   migrator.MigrationUp,
		}
		err := fsRepo.CreateMigrationFile(migration)
		if err != nil {
			t.Error(err)
		}

		filenames, err := fsRepo.GetMigrationFilePathsByType(migrator.MigrationUp)
		if err != nil {
			t.Error(err)
		}

		if !testutils.Contains(filenames, filepath.Join(testDir, "000100_create_users_table_up."+extention)) {
			t.Errorf("expected %s to be in filenames", migration.Name)
		}
	})

	t.Run("create down migration file", func(t *testing.T) {
		migration := &migrator.Migration{
			Number: 100,
			Name:   "create_users_table",
			Type:   migrator.MigrationDown,
		}
		err := fsRepo.CreateMigrationFile(migration)
		if err != nil {
			t.Error(err)
		}

		filenames, err := fsRepo.GetMigrationFilePathsByType(migrator.MigrationDown)
		if err != nil {
			t.Error(err)
		}

		if !testutils.Contains(filenames, filepath.Join(testDir, "000100_create_users_table_down."+extention)) {
			t.Errorf("expected %s to be in filenames", migration.Name)
		}
	})
}
