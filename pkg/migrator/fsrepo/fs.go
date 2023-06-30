package fsrepo

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adharshmk96/stk/pkg/migrator"
)

type fileSystem struct {
	workDir string
	ext     string
}

func NewFSRepo(workDir, ext string) migrator.FileRepo {
	return &fileSystem{
		workDir: workDir,
		ext:     ext,
	}
}

func (f *fileSystem) LoadMigrationsFromFile(migrationType migrator.MigrationType) ([]*migrator.Migration, error) {
	if err := f.OpenDirectory(); err != nil {
		return nil, err
	}
	filePaths, err := f.GetMigrationFilePathsByType(migrationType)
	if err != nil {
		return nil, err
	}

	migrations, err := parseMigrationsFromFilePaths(filePaths)
	if err != nil {
		return nil, err
	}

	if migrationType == migrator.MigrationDown {
		sortDescMigrations(migrations)
	} else {
		sortAscMigrations(migrations)
	}

	return migrations, nil
}

func (f *fileSystem) OpenDirectory() error {
	return mkPathIfNotExists(f.workDir)
}

func (f *fileSystem) GetMigrationFilePathsByType(migrationType migrator.MigrationType) ([]string, error) {
	entries, err := os.ReadDir(f.workDir)
	if err != nil {
		return nil, err
	}

	migrationFilePaths := make([]string, 0, len(entries))
	for _, entry := range entries {
		// check the string ends with "down"
		if !entry.IsDir() {
			filenameWithoutExt := fileNameWithoutExtension(entry.Name())
			if strings.HasSuffix(filenameWithoutExt, string(migrationType)) {
				fullPath := filepath.Join(f.workDir, entry.Name())
				migrationFilePaths = append(migrationFilePaths, fullPath)
			} else {
				continue
			}
		}
	}

	return migrationFilePaths, nil
}

func (f *fileSystem) CreateMigrationFile(migration *migrator.Migration) error {
	fileName := migrator.MigrationToFilename(migration)
	filePath := filepath.Join(f.workDir, addExtension(fileName, f.ext))
	migration.Path = filePath

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func (f *fileSystem) WriteMigrationToFile(migration *migrator.Migration) error {
	err := os.WriteFile(migration.Path, []byte(migration.Query), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (f *fileSystem) LoadMigrationQuery(migration *migrator.Migration) error {
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return err
	}
	migration.Query = string(content)
	return nil
}

func mkPathIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func parseMigrationsFromFilePaths(filePaths []string) ([]*migrator.Migration, error) {
	migrations := make([]*migrator.Migration, 0, len(filePaths))
	for _, filePath := range filePaths {
		nameWithExt := filepath.Base(filePath)
		ext := filepath.Ext(nameWithExt)
		nameWithoutExt := strings.TrimSuffix(nameWithExt, ext)
		migration, err := migrator.ParseMigrationFromString(nameWithoutExt)
		if err != nil {
			return nil, err
		}
		migration.Path = filePath
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func sortAscMigrations(migrations []*migrator.Migration) {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Number < migrations[j].Number
	})
}

func sortDescMigrations(migrations []*migrator.Migration) {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Number > migrations[j].Number
	})
}

func addExtension(filename, extension string) string {
	// Remove existing extension if there is one
	currentExtension := filepath.Ext(filename)
	if currentExtension != "" {
		filename = strings.TrimSuffix(filename, currentExtension)
	}

	// Ensure the new extension starts with a dot
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	return filename + extension
}
