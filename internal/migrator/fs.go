package migrator

import (
	"os"
	"path/filepath"
	"strings"
)

func GetFilenamesWithoutExtension(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			filenameWithoutExt := fileNameWithoutExtension(entry.Name())
			filenames = append(filenames, filenameWithoutExt)
		}
	}

	return filenames, nil
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
