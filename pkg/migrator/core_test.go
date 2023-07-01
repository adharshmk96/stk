package migrator

import (
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
				actual, _ := ParseMigrationFromString(c.fileName)

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
				mig, err := ParseMigrationFromString(c.fileName)
				t.Log(mig)

				assert.Equal(t, c.expected, err)
			})
		}
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
				actual := MigrationToFilename(c.migration)

				assert.Equal(t, c.expected, actual)
			})
		}
	})

}
