package migrator_test

import (
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
				fileName: "2_up",
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
