package testutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func removeFile(t *testing.T, filePath string) func() {
	return func() {
		err := os.Remove(filePath)
		assert.NoError(t, err)
	}
}

func CreateTempFile(t *testing.T, content string) (string, func()) {
	file, err := os.CreateTemp("", "test")
	assert.NoError(t, err)

	defer file.Close()

	_, err = file.WriteString(content)
	assert.NoError(t, err)

	return file.Name(), removeFile(t, file.Name())
}

func removeDir(t *testing.T, dir string) func() {
	return func() {
		err := os.RemoveAll(dir)
		assert.NoError(t, err)
	}
}

func CreateTempDirectory(t *testing.T) (string, func()) {
	dir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	return dir, removeDir(t, dir)
}

func GetFileContent(t *testing.T, filePath string) string {
	t.Helper()
	file, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	return string(file)
}

func GetNumberOfFilesInFolder(t *testing.T, folder string) int {
	t.Helper()
	files, err := os.ReadDir(folder)
	assert.NoError(t, err)
	return len(files)
}
