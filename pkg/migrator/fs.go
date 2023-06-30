package migrator

import (
	"os"
	"path/filepath"
	"strings"
)

func openDirectory(root string, database Database) string {
	var subDirectory string
	switch database {
	case PostgresDB:
		subDirectory = "postgres"
	case MySQLDB:
		subDirectory = "mysql"
	case SQLiteDB:
		subDirectory = "sqlite"
	default:
		subDirectory = "sqlite"
	}

	directory := filepath.Join(root, subDirectory)
	mkPathIfNotExists(directory)

	return directory
}

// TODO: return full path instead and let parser handle the rest.
func getMigrationFileGroup(dir string, migrationType MigrationType) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	migrationFiles := make([]string, 0, len(entries))
	for _, entry := range entries {
		// check the string ends with "down"
		if !entry.IsDir() {
			filenameWithoutExt := fileNameWithoutExtension(entry.Name())
			if strings.HasSuffix(filenameWithoutExt, string(migrationType)) {
				migrationFiles = append(migrationFiles, filenameWithoutExt)
			}
		}
	}

	return migrationFiles, nil
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func mkPathIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func createFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}

func createMigrationFile(dir string, migrationFileName string) error {
	path := filepath.Join(dir, migrationFileName)
	return createFile(path)
}
