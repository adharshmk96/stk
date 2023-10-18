package sqlmigrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNextMigrations(t *testing.T) {
	t.Run("generates next migrations", func(t *testing.T) {
		lastMigration := Migration{
			Number: 1,
			Name:   "create_users_table",
		}

		migrations := []Migration{
			{
				Number: 2,
			},
			{
				Number: 3,
			},
		}

		nextMigrations := GenerateNextMigrations(lastMigration.Number, "", len(migrations))

		for i, migration := range nextMigrations {
			assert.Equal(t, migrations[i].Number, migration.Number)
		}

	})
}
