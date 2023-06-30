package migrator_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/adharshmk96/stk/mocks"
	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_findUpMigrationsToApply(t *testing.T) {

	var migrations []*migrator.Migration
	for i := 1; i <= 10; i++ {
		migrations = append(migrations, &migrator.Migration{
			Number: i,
			Name:   "test" + strconv.Itoa(i),
			Path:   "test" + strconv.Itoa(i),
			Query:  "test" + strconv.Itoa(i),
			Type:   migrator.MigrationUp,
		})
	}

	type args struct {
		lastMigration   *migrator.Migration
		migrations      []*migrator.Migration
		numberToMigrate int
	}
	tests := []struct {
		name string
		args args
		want []*migrator.Migration
	}{
		{
			name: "find up migrations after num 5 to apply",
			args: args{
				lastMigration: &migrator.Migration{
					Number: 5,
					Name:   "test",
					Path:   "test",
					Query:  "test",
					Type:   migrator.MigrationUp,
				},
				migrations:      migrations,
				numberToMigrate: 2,
			},
			want: migrations[5 : 5+2],
		},
		{
			name: "find up migrations from num 5 to apply",
			args: args{
				lastMigration: &migrator.Migration{
					Number: 5,
					Name:   "test",
					Path:   "test",
					Query:  "test",
					Type:   migrator.MigrationDown,
				},
				migrations:      migrations,
				numberToMigrate: 2,
			},
			want: migrations[4 : 4+2],
		},
		{
			name: "apply 20 migrations from num 5",
			args: args{
				lastMigration: &migrator.Migration{
					Number: 5,
					Name:   "test",
					Path:   "test",
					Query:  "test",
					Type:   migrator.MigrationDown,
				},
				migrations:      migrations,
				numberToMigrate: 20,
			},
			want: migrations[4:],
		},
		{
			name: "apply 20 migrations after num 5",
			args: args{
				lastMigration: &migrator.Migration{
					Number: 5,
					Name:   "test",
					Path:   "test",
					Query:  "test",
					Type:   migrator.MigrationUp,
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
			name: "return first n migrations if last migration is nil",
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
			want: []*migrator.Migration{},
		},
		{
			name: "last migration is nil and num to migrate is 0",
			args: args{
				lastMigration:   nil,
				migrations:      migrations,
				numberToMigrate: 0,
			},
			want: migrations,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := migrator.FindUpMigrationsToApply(tt.args.lastMigration, tt.args.migrations, tt.args.numberToMigrate)

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

func TestApplyMigration(t *testing.T) {

	t.Run("Test dry run", func(t *testing.T) {
		mockDBRepo := mocks.NewDatabaseRepo(t)
		mockMigration := &migrator.Migration{
			Number: 001,
			Name:   "test_migration",
			Type:   migrator.MigrationUp,
		}
		mockConfig := &migrator.MigratorConfig{
			DryRun: true,
			DBRepo: mockDBRepo,
		}

		err := migrator.ApplyMigration(mockConfig, mockMigration)
		assert.NoError(t, err)

		mockDBRepo.AssertNotCalled(t, "ApplyMigration")
	})

	t.Run("Test apply migration success", func(t *testing.T) {
		mockDBRepo := mocks.NewDatabaseRepo(t)
		mockMigration := &migrator.Migration{
			Number: 001,
			Name:   "test_migration",
			Type:   migrator.MigrationUp,
		}

		mockConfig := &migrator.MigratorConfig{
			DryRun: false,
			DBRepo: mockDBRepo,
		}

		mockDBRepo.On("ApplyMigration", mockMigration).Return(nil)

		err := migrator.ApplyMigration(mockConfig, mockMigration)
		assert.NoError(t, err)
		mockDBRepo.AssertExpectations(t)

		mockDBRepo.AssertNumberOfCalls(t, "ApplyMigration", 1)
	})

	t.Run("Test apply migration failure", func(t *testing.T) {
		mockDBRepo := mocks.NewDatabaseRepo(t)
		mockMigration := &migrator.Migration{
			Number: 001,
			Name:   "test_migration",
			Type:   migrator.MigrationUp,
		}
		mockConfig := &migrator.MigratorConfig{
			DryRun: false,
			DBRepo: mockDBRepo,
		}

		mockDBRepo.On("ApplyMigration", mock.Anything).Return(errors.New("db error"))

		err := migrator.ApplyMigration(mockConfig, mockMigration)
		assert.Error(t, err)
		mockDBRepo.AssertExpectations(t)
		mockDBRepo.AssertNumberOfCalls(t, "ApplyMigration", 1)
	})
}

func TestMigrateUp(t *testing.T) {
	t.Run("migrate up on empty database and empty file repo", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		dbRepo := mocks.NewDatabaseRepo(t)
		mockConfig := &migrator.MigratorConfig{
			FSRepo:       fsRepo,
			DBRepo:       dbRepo,
			NumToMigrate: 5,
			DryRun:       false,
		}

		var migrations []*migrator.Migration
		// for i := 0; i < 5; i++ {
		// 	migrations = append(migrations, &migrator.Migration{
		// 		Number: i,
		// 		Name:   fmt.Sprintf("test_%d", i),
		// 		Type:   migrator.MigrationUp,
		// 	})
		// }

		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)
		// dbRepo.On("ApplyMigration").Return(migrations, nil)

		migrations, err := migrator.MigrateUp(mockConfig)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(migrations))

		dbRepo.AssertNotCalled(t, "ApplyMigration")

	})

	t.Run("migrate up on non-empty database and empty file repo", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		dbRepo := mocks.NewDatabaseRepo(t)
		mockConfig := &migrator.MigratorConfig{
			FSRepo:       fsRepo,
			DBRepo:       dbRepo,
			NumToMigrate: 5,
			DryRun:       false,
		}

		var migrations []*migrator.Migration
		for i := 0; i < 5; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		dbRepo.On("LoadLastAppliedMigration").Return(migrations[2], nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return([]*migrator.Migration{}, nil)

		migrations, err := migrator.MigrateUp(mockConfig)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(migrations))

		dbRepo.AssertNotCalled(t, "ApplyMigration")
		dbRepo.AssertNotCalled(t, "LoadMigrationQuery")

	})

	t.Run("migrate up with dry run", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		dbRepo := mocks.NewDatabaseRepo(t)
		mockConfig := &migrator.MigratorConfig{
			FSRepo:       fsRepo,
			DBRepo:       dbRepo,
			NumToMigrate: 5,
			DryRun:       true,
		}

		var migrations []*migrator.Migration
		for i := 0; i < 5; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
		fsRepo.On("LoadMigrationQuery", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)

		migrations, err := migrator.MigrateUp(mockConfig)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(migrations))

		dbRepo.AssertNotCalled(t, "ApplyMigration")
	})

	t.Run("migrate up on empty database and non-empty file repo", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		dbRepo := mocks.NewDatabaseRepo(t)
		mockConfig := &migrator.MigratorConfig{
			FSRepo:       fsRepo,
			DBRepo:       dbRepo,
			NumToMigrate: 5,
			DryRun:       false,
		}

		var migrations []*migrator.Migration
		for i := 1; i <= 5; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)
		fsRepo.On("LoadMigrationQuery", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
		dbRepo.On("ApplyMigration", mock.AnythingOfType("*migrator.Migration")).Return(nil)

		migrations, err := migrator.MigrateUp(mockConfig)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(migrations))

		fsRepo.AssertExpectations(t)
		dbRepo.AssertExpectations(t)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationsFromFile", 1)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationQuery", 5)
		dbRepo.AssertNumberOfCalls(t, "LoadLastAppliedMigration", 1)
		dbRepo.AssertNumberOfCalls(t, "ApplyMigration", 5)
	})

	t.Run("migrate up on non-empty database and non-empty file repo", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		dbRepo := mocks.NewDatabaseRepo(t)
		mockConfig := &migrator.MigratorConfig{
			FSRepo:       fsRepo,
			DBRepo:       dbRepo,
			NumToMigrate: 5,
			DryRun:       false,
		}

		var migrations []*migrator.Migration
		for i := 1; i <= 5; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		dbRepo.On("LoadLastAppliedMigration").Return(migrations[2], nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)
		fsRepo.On("LoadMigrationQuery", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		dbRepo.On("ApplyMigration", mock.AnythingOfType("*migrator.Migration")).Return(nil)

		migrations, err := migrator.MigrateUp(mockConfig)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(migrations))

		fsRepo.AssertExpectations(t)
		dbRepo.AssertExpectations(t)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationsFromFile", 1)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationQuery", 2)
		dbRepo.AssertNumberOfCalls(t, "LoadLastAppliedMigration", 1)
		dbRepo.AssertNumberOfCalls(t, "ApplyMigration", 2)
	})

	t.Run("migrate up with 0 num to migrate", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		dbRepo := mocks.NewDatabaseRepo(t)
		mockConfig := &migrator.MigratorConfig{
			FSRepo:       fsRepo,
			DBRepo:       dbRepo,
			NumToMigrate: 0,
			DryRun:       false,
		}

		var migrations []*migrator.Migration
		for i := 1; i <= 5; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)
		fsRepo.On("LoadMigrationQuery", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		dbRepo.On("ApplyMigration", mock.AnythingOfType("*migrator.Migration")).Return(nil)

		migrations, err := migrator.MigrateUp(mockConfig)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(migrations))

		fsRepo.AssertExpectations(t)
		dbRepo.AssertExpectations(t)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationsFromFile", 1)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationQuery", 5)
		dbRepo.AssertNumberOfCalls(t, "LoadLastAppliedMigration", 1)
		dbRepo.AssertNumberOfCalls(t, "ApplyMigration", 5)
	})
}
