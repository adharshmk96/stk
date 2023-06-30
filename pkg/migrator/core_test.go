package migrator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStringToMigration(t *testing.T) {
	t.Run("parse string to migration", func(t *testing.T) {
		tc := []struct {
			fileName string
			expected *Migration
		}{
			{
				fileName: "1000002_up",
				expected: &Migration{
					Number: 1000002,
					Name:   "",
					Type:   MigrationUp,
				},
			},
			{
				fileName: "1000002__down",
				expected: &Migration{
					Number: 1000002,
					Name:   "",
					Type:   MigrationDown,
				},
			},
			{
				fileName: "2__up",
				expected: &Migration{
					Number: 2,
					Name:   "",
					Type:   MigrationUp,
				},
			},
			{
				fileName: "2_down",
				expected: &Migration{
					Number: 2,
					Name:   "",
					Type:   MigrationDown,
				},
			},
			{
				fileName: "1_create_users_table_up",
				expected: &Migration{
					Number: 1,
					Name:   "create_users_table",
					Type:   MigrationUp,
				},
			},
			{
				fileName: "1_create_users_table_down",
				expected: &Migration{
					Number: 1,
					Name:   "create_users_table",
					Type:   MigrationDown,
				},
			},
			{
				fileName: "000002_create_posts_table_up",
				expected: &Migration{
					Number: 2,
					Name:   "create_posts_table",
					Type:   MigrationUp,
				},
			},
			{
				fileName: "000002_create_posts_table_down",
				expected: &Migration{
					Number: 2,
					Name:   "create_posts_table",
					Type:   MigrationDown,
				},
			},
			{
				fileName: "1000002_create_posts_table_up",
				expected: &Migration{
					Number: 1000002,
					Name:   "create_posts_table",
					Type:   MigrationUp,
				},
			},
			{
				fileName: "1000002_create_posts_table_down",
				expected: &Migration{
					Number: 1000002,
					Name:   "create_posts_table",
					Type:   MigrationDown,
				},
			},
		}

		for _, c := range tc {
			t.Run(c.fileName, func(t *testing.T) {
				actual, _ := parseMigrationFromString(c.fileName)

				assert.Equal(t, c.expected.Number, actual.Number)
				assert.Equal(t, c.expected.Name, actual.Name)
				assert.Equal(t, c.expected.Type, actual.Type)

			})

		}
	})

	t.Run("parse string to migration error", func(t *testing.T) {
		tc := []struct {
			fileName string
			expected error
		}{
			{
				fileName: "1000002",
				expected: ErrInvalidFormat,
			},
			{
				fileName: "1000002_",
				expected: ErrInvalidFormat,
			},
			{
				fileName: "up",
				expected: ErrInvalidFormat,
			},
			{
				fileName: "a_b",
				expected: ErrInvalidFormat,
			},
		}

		for _, c := range tc {
			t.Run(c.fileName, func(t *testing.T) {
				mig, err := parseMigrationFromString(c.fileName)
				t.Log(mig)

				assert.Equal(t, c.expected, err)
			})
		}
	})
}

func TestParseMigrationsFromFilenames(t *testing.T) {
	t.Run("parse migrations from filenames", func(t *testing.T) {
		fileNames := []string{
			"1000002_up",
			"1000002_down",
			"1000001_create_users_table_up",
			"1000001_create_users_table_down",
		}

		expected := []*Migration{
			{
				Number: 1000002,
				Name:   "",
				Type:   MigrationUp,
			},
			{
				Number: 1000002,
				Name:   "",
				Type:   MigrationDown,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   MigrationUp,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   MigrationDown,
			},
		}

		actual, _ := parseMigrationsFromFilenames(fileNames)

		assert.Equal(t, expected, actual)
	})
}

func TestSortMigrations(t *testing.T) {
	t.Run("sorts migration objects", func(t *testing.T) {
		migrations := []*Migration{
			{
				Number: 1000002,
				Name:   "",
				Type:   MigrationUp,
			},
			{
				Number: 1000002,
				Name:   "",
				Type:   MigrationDown,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   MigrationUp,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   MigrationDown,
			},
			{
				Number: 999999,
				Name:   "create_posts_table",
				Type:   MigrationUp,
			},
			{
				Number: 1,
				Name:   "create_posts_table",
				Type:   MigrationDown,
			},
		}

		expected := []*Migration{
			{
				Number: 1,
				Name:   "create_posts_table",
				Type:   MigrationDown,
			},
			{
				Number: 999999,
				Name:   "create_posts_table",
				Type:   MigrationUp,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   MigrationUp,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   MigrationDown,
			},
			{
				Number: 1000002,
				Name:   "",
				Type:   MigrationUp,
			},
			{
				Number: 1000002,
				Name:   "",
				Type:   MigrationDown,
			},
		}

		sortMigrations(migrations)

		assert.Equal(t, expected, migrations)
	})
}

func TestGenerateNextMigrations(t *testing.T) {
	t.Run("generate next migrations after n", func(t *testing.T) {
		tc := []struct {
			starting int
			totalNum int
		}{
			{
				starting: 1,
				totalNum: 10,
			},
			{
				starting: 999999,
				totalNum: 10,
			},
			{
				starting: 1000000,
				totalNum: 10,
			},
		}

		for _, c := range tc {
			t.Run(fmt.Sprintf("starting at %d", c.starting), func(t *testing.T) {
				expected := make([]*Migration, 0)

				for i := c.starting + 1; i <= c.starting+c.totalNum; i++ {
					expected = append(expected, &Migration{
						Number: int(i),
						Name:   "create_users_table",
						Type:   MigrationUp,
					})
					expected = append(expected, &Migration{
						Number: int(i),
						Name:   "create_users_table",
						Type:   MigrationDown,
					})
				}

				actual := generateNextMigrations(c.starting, "create_users_table", c.totalNum)

				for i := range actual {
					assert.Equal(t, expected[i].Number, actual[i].Number)
					assert.Equal(t, expected[i].Name, actual[i].Name)
					assert.Equal(t, expected[i].Type, actual[i].Type)
				}
			})
		}

	})
}

func TestMigrationToFilename(t *testing.T) {
	t.Run("converts migration to filename", func(t *testing.T) {
		tc := []struct {
			migration *Migration
			expected  string
		}{
			{
				migration: &Migration{
					Number: 1000001,
					Name:   "create_users_table",
					Type:   MigrationUp,
				},
				expected: "1000001_create_users_table_up",
			},
			{
				migration: &Migration{
					Number: 1000001,
					Name:   "create_users_table",
					Type:   MigrationDown,
				},
				expected: "1000001_create_users_table_down",
			},
			{
				migration: &Migration{
					Number: 1000001,
					Name:   "",
					Type:   MigrationUp,
				},
				expected: "1000001__up",
			},
			{
				migration: &Migration{
					Number: 1000001,
					Name:   "",
					Type:   MigrationDown,
				},
				expected: "1000001__down",
			},
		}

		for _, c := range tc {
			t.Run(c.expected, func(t *testing.T) {
				actual := migrationToFilename(c.migration)

				assert.Equal(t, c.expected, actual)
			})
		}
	})

}

func TestParseMigrationsFromFilePaths(t *testing.T) {
	t.Run("parses migrations from file paths", func(t *testing.T) {
		var filepaths []string
		expected := make([]*Migration, 0)

		for i := 1; i <= 10; i++ {
			name := fmt.Sprintf("create_users_table%d", i)

			up := fmt.Sprintf("%06d_%s_up.sql", i, name)
			down := fmt.Sprintf("%06d_%s_down.sql", i, name)
			filepaths = append(filepaths, up)
			filepaths = append(filepaths, down)

			expected = append(expected, &Migration{
				Number: i,
				Name:   name,
				Type:   MigrationUp,
				Path:   up,
			})
			expected = append(expected, &Migration{
				Number: i,
				Name:   name,
				Type:   MigrationDown,
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
