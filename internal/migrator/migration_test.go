package migrator_test

import (
	"fmt"
	"testing"

	"github.com/adharshmk96/stk/internal/migrator"
	"github.com/stretchr/testify/assert"
)

func TestParseStringToMigration(t *testing.T) {
	t.Run("parse string to migration", func(t *testing.T) {
		tc := []struct {
			fileName string
			expected *migrator.Migration
		}{
			{
				fileName: "1000002_up",
				expected: &migrator.Migration{
					Number: 1000002,
					Name:   "",
					Type:   migrator.MigrationUp,
				},
			},
			{
				fileName: "1000002__down",
				expected: &migrator.Migration{
					Number: 1000002,
					Name:   "",
					Type:   migrator.MigrationDown,
				},
			},
			{
				fileName: "2__up",
				expected: &migrator.Migration{
					Number: 2,
					Name:   "",
					Type:   migrator.MigrationUp,
				},
			},
			{
				fileName: "2_down",
				expected: &migrator.Migration{
					Number: 2,
					Name:   "",
					Type:   migrator.MigrationDown,
				},
			},
			{
				fileName: "1_create_users_table_up",
				expected: &migrator.Migration{
					Number: 1,
					Name:   "create_users_table",
					Type:   migrator.MigrationUp,
				},
			},
			{
				fileName: "1_create_users_table_down",
				expected: &migrator.Migration{
					Number: 1,
					Name:   "create_users_table",
					Type:   migrator.MigrationDown,
				},
			},
			{
				fileName: "000002_create_posts_table_up",
				expected: &migrator.Migration{
					Number: 2,
					Name:   "create_posts_table",
					Type:   migrator.MigrationUp,
				},
			},
			{
				fileName: "000002_create_posts_table_down",
				expected: &migrator.Migration{
					Number: 2,
					Name:   "create_posts_table",
					Type:   migrator.MigrationDown,
				},
			},
			{
				fileName: "1000002_create_posts_table_up",
				expected: &migrator.Migration{
					Number: 1000002,
					Name:   "create_posts_table",
					Type:   migrator.MigrationUp,
				},
			},
			{
				fileName: "1000002_create_posts_table_down",
				expected: &migrator.Migration{
					Number: 1000002,
					Name:   "create_posts_table",
					Type:   migrator.MigrationDown,
				},
			},
		}

		for _, c := range tc {
			t.Run(c.fileName, func(t *testing.T) {
				actual, _ := migrator.ParseMigration(c.fileName)

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
				expected: migrator.ErrInvalidFormat,
			},
			{
				fileName: "1000002_",
				expected: migrator.ErrInvalidFormat,
			},
			{
				fileName: "up",
				expected: migrator.ErrInvalidFormat,
			},
			{
				fileName: "a_b",
				expected: migrator.ErrInvalidFormat,
			},
		}

		for _, c := range tc {
			t.Run(c.fileName, func(t *testing.T) {
				mig, err := migrator.ParseMigration(c.fileName)
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

		expected := []*migrator.Migration{
			{
				Number: 1000002,
				Name:   "",
				Type:   migrator.MigrationUp,
			},
			{
				Number: 1000002,
				Name:   "",
				Type:   migrator.MigrationDown,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   migrator.MigrationUp,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   migrator.MigrationDown,
			},
		}

		actual, _ := migrator.ParseMigrationsFromFilenames(fileNames)

		assert.Equal(t, expected, actual)
	})
}

func TestSortMigrations(t *testing.T) {
	t.Run("sorts migration objects", func(t *testing.T) {
		migrations := []*migrator.Migration{
			{
				Number: 1000002,
				Name:   "",
				Type:   migrator.MigrationUp,
			},
			{
				Number: 1000002,
				Name:   "",
				Type:   migrator.MigrationDown,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   migrator.MigrationUp,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   migrator.MigrationDown,
			},
			{
				Number: 999999,
				Name:   "create_posts_table",
				Type:   migrator.MigrationUp,
			},
			{
				Number: 1,
				Name:   "create_posts_table",
				Type:   migrator.MigrationDown,
			},
		}

		expected := []*migrator.Migration{
			{
				Number: 1,
				Name:   "create_posts_table",
				Type:   migrator.MigrationDown,
			},
			{
				Number: 999999,
				Name:   "create_posts_table",
				Type:   migrator.MigrationUp,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   migrator.MigrationUp,
			},
			{
				Number: 1000001,
				Name:   "create_users_table",
				Type:   migrator.MigrationDown,
			},
			{
				Number: 1000002,
				Name:   "",
				Type:   migrator.MigrationUp,
			},
			{
				Number: 1000002,
				Name:   "",
				Type:   migrator.MigrationDown,
			},
		}

		migrator.SortMigrations(migrations)

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
				expected := make([]*migrator.Migration, 0)

				for i := c.starting + 1; i <= c.starting+c.totalNum; i++ {
					expected = append(expected, &migrator.Migration{
						Number: int(i),
						Name:   "create_users_table",
						Type:   migrator.MigrationUp,
					})
					expected = append(expected, &migrator.Migration{
						Number: int(i),
						Name:   "create_users_table",
						Type:   migrator.MigrationDown,
					})
				}

				actual := migrator.GenerateNextMigrations(c.starting, "create_users_table", c.totalNum)

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
			migration *migrator.Migration
			expected  string
		}{
			{
				migration: &migrator.Migration{
					Number: 1000001,
					Name:   "create_users_table",
					Type:   migrator.MigrationUp,
				},
				expected: "1000001_create_users_table_up",
			},
			{
				migration: &migrator.Migration{
					Number: 1000001,
					Name:   "create_users_table",
					Type:   migrator.MigrationDown,
				},
				expected: "1000001_create_users_table_down",
			},
			{
				migration: &migrator.Migration{
					Number: 1000001,
					Name:   "",
					Type:   migrator.MigrationUp,
				},
				expected: "1000001__up",
			},
			{
				migration: &migrator.Migration{
					Number: 1000001,
					Name:   "",
					Type:   migrator.MigrationDown,
				},
				expected: "1000001__down",
			},
		}

		for _, c := range tc {
			t.Run(c.expected, func(t *testing.T) {
				actual := migrator.MigrationToFilename(c.migration)

				assert.Equal(t, c.expected, actual)
			})
		}
	})
}
