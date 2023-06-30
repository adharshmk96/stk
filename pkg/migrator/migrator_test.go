package migrator

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateUp(t *testing.T) {
	t.Run("migrate up", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_findUpMigrationsToApply(t *testing.T) {

	var migrations []*Migration
	for i := 1; i <= 10; i++ {
		migrations = append(migrations, &Migration{
			Number: i,
			Name:   "test" + strconv.Itoa(i),
			Path:   "test" + strconv.Itoa(i),
			Query:  "test" + strconv.Itoa(i),
			Type:   MigrationUp,
		})
	}

	type args struct {
		lastMigration   *Migration
		migrations      []*Migration
		numberToMigrate int
	}
	tests := []struct {
		name string
		args args
		want []*Migration
	}{
		{
			name: "find up migrations after num 5 to apply",
			args: args{
				lastMigration: &Migration{
					Number: 5,
					Name:   "test",
					Path:   "test",
					Query:  "test",
					Type:   MigrationUp,
				},
				migrations:      migrations,
				numberToMigrate: 2,
			},
			want: migrations[5 : 5+2],
		},
		{
			name: "find up migrations from num 5 to apply",
			args: args{
				lastMigration: &Migration{
					Number: 5,
					Name:   "test",
					Path:   "test",
					Query:  "test",
					Type:   MigrationDown,
				},
				migrations:      migrations,
				numberToMigrate: 2,
			},
			want: migrations[4 : 4+2],
		},
		{
			name: "apply 20 migrations from num 5",
			args: args{
				lastMigration: &Migration{
					Number: 5,
					Name:   "test",
					Path:   "test",
					Query:  "test",
					Type:   MigrationDown,
				},
				migrations:      migrations,
				numberToMigrate: 20,
			},
			want: migrations[4:],
		},
		{
			name: "apply 20 migrations after num 5",
			args: args{
				lastMigration: &Migration{
					Number: 5,
					Name:   "test",
					Path:   "test",
					Query:  "test",
					Type:   MigrationUp,
				},
				migrations:      migrations,
				numberToMigrate: 20,
			},
			want: migrations[5:],
		},
		{
			name: "return all migrations if last migration is nil",
			args: args{
				lastMigration:   nil,
				migrations:      migrations,
				numberToMigrate: 20,
			},
			want: migrations,
		},
		{
			name: "return all migrations if last migration is nil",
			args: args{
				lastMigration:   nil,
				migrations:      migrations,
				numberToMigrate: 2,
			},
			want: migrations[:2],
		},
		{
			name: "last migration is the last migration in the list",
			args: args{
				lastMigration:   migrations[len(migrations)-1],
				migrations:      migrations,
				numberToMigrate: 2,
			},
			want: []*Migration{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findUpMigrationsToApply(tt.args.lastMigration, tt.args.migrations, tt.args.numberToMigrate)

			if tt.name == "last migration is the last migration in the list" {
				assert.Equal(t, 0, len(got))
				return
			}

			assert.Equal(t, len(tt.want), len(got))

			for idx, migration := range tt.want {
				assert.Equal(t, migration.Number, got[idx].Number)
				assert.Equal(t, migration.Name, got[idx].Name)
				assert.Equal(t, migration.Path, got[idx].Path)
				assert.Equal(t, migration.Query, got[idx].Query)
			}
		})
	}
}
