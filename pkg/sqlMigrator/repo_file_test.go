package sqlmigrator

import (
	"testing"

	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

var TEST_FILE_CONTENT = `1_create_users_table_up
2_create_posts_table_up
3_create_comments_table_up
4_create_likes_table_down
5_create_followers_table_down
6_create_messages_table_down
`

func TestReadLines(t *testing.T) {
	t.Run("reads a file and returns all lines", func(t *testing.T) {
		filePath, removeTempFile := testutils.CreateTempFile(t, TEST_FILE_CONTENT)
		lines, err := ReadLines(filePath)

		defer removeTempFile()

		expected := []string{
			"1_create_users_table_up",
			"2_create_posts_table_up",
			"3_create_comments_table_up",
			"4_create_likes_table_down",
			"5_create_followers_table_down",
			"6_create_messages_table_down",
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, lines)
	})

	t.Run("returns an error if file does not exist", func(t *testing.T) {
		_, err := ReadLines("non-existent-file.txt")

		assert.Error(t, err)
	})
}

func TestReadContent(t *testing.T) {

}
