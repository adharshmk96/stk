package migrator

import (
	"os"
	"path/filepath"
	"strings"
)

func GetMigrationFileGroup(dir string, migrationType MigrationType) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	migrationFiles := make([]string, 0, len(entries))
	for _, entry := range entries {
		// check the string ends with "down"
		if !entry.IsDir() {
			filenameWithoutExt := fileNameWithoutExtension(entry.Name())
			if !strings.HasSuffix(filenameWithoutExt, string(migrationType)) {
				migrationFiles = append(migrationFiles, entry.Name())
			}
		}
	}

	return migrationFiles, nil
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func MkPathIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0700)
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

func CreateMigrationFile(dir string, migrationFileName string) error {
	path := filepath.Join(dir, migrationFileName)
	return createFile(path)
}
