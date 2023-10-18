package sqlmigrator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TEST_FILE_CONTENT = `1_create_users_table
2_create_posts_table
3_create_comments_table
`

func createTempFilePath(t *testing.T) string {
	file, err := os.CreateTemp("", "test")
	assert.NoError(t, err)

	return file.Name()
}

func createTempFile(t *testing.T, content string) string {
	filePath := createTempFilePath(t)

	file, err := os.Create(filePath)
	assert.NoError(t, err)

	defer file.Close()

	_, err = file.WriteString(content)
	assert.NoError(t, err)

	return filePath
}

func removeTempFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	assert.NoError(t, err)
}

func TestReadLastLine(t *testing.T) {
	t.Run("reads the last line of a file", func(t *testing.T) {
		filePath := createTempFile(t, TEST_FILE_CONTENT)
		defer removeTempFile(t, filePath)

		line, err := readLastLine(filePath)
		assert.NoError(t, err)

		assert.Equal(t, "3_create_comments_table", line)
	})
}
