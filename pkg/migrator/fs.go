package migrator

import (
	"fmt"
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

	fmt.Println(directory)

	return directory
}

func getMigrationFilePathsByGroup(dir string, migrationType MigrationType) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	migrationFilePaths := make([]string, 0, len(entries))
	for _, entry := range entries {
		// check the string ends with "down"
		if !entry.IsDir() {
			filenameWithoutExt := fileNameWithoutExtension(entry.Name())
			if strings.HasSuffix(filenameWithoutExt, string(migrationType)) {
				fullPath := filepath.Join(dir, entry.Name())
				migrationFilePaths = append(migrationFilePaths, fullPath)
			} else {
				continue
			}
		}
	}

	return migrationFilePaths, nil
}

func getMigrationFileNamesByGroup(dir string, migrationType MigrationType) ([]string, error) {
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

func writeToFile(path string, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func readFileContents(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
