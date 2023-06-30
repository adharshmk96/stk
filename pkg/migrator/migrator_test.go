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

func TestMigrateDown(t *testing.T) {
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

		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return(migrations, nil)
		// dbRepo.On("ApplyMigration").Return(migrations, nil)

		migrations, err := migrator.MigrateDown(mockConfig)
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
		for i := 5; i > 0; i-- {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		dbRepo.On("LoadLastAppliedMigration").Return(migrations[2], nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return([]*migrator.Migration{}, nil)

		migrations, err := migrator.MigrateDown(mockConfig)
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
		for i := 5; i > 0; i-- {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
		fsRepo.On("LoadMigrationQuery", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return(migrations, nil)

		migrations, err := migrator.MigrateDown(mockConfig)
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
		for i := 5; i > 0; i-- {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return(migrations, nil)
		fsRepo.On("LoadMigrationQuery", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
		dbRepo.On("ApplyMigration", mock.AnythingOfType("*migrator.Migration")).Return(nil)

		migrations, err := migrator.MigrateDown(mockConfig)
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
		for i := 5; i > 0; i-- {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		dbRepo.On("LoadLastAppliedMigration").Return(migrations[2], nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return(migrations, nil)
		fsRepo.On("LoadMigrationQuery", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		dbRepo.On("ApplyMigration", mock.AnythingOfType("*migrator.Migration")).Return(nil)

		migrations, err := migrator.MigrateDown(mockConfig)
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
		for i := 5; i > 0; i-- {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return(migrations, nil)
		fsRepo.On("LoadMigrationQuery", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		dbRepo.On("ApplyMigration", mock.AnythingOfType("*migrator.Migration")).Return(nil)

		migrations, err := migrator.MigrateDown(mockConfig)
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

func TestCalculateDownMigration(t *testing.T) {
	t.Run("calculate down migration with last up migration", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		var expected []*migrator.Migration
		for i := 5; i > 0; i-- {
			expected = append(expected, migrations[i-1])
		}

		lastUpMigration := migrations[4]

		downMigrations := migrator.CalculateDownMigrationsToApply(lastUpMigration, migrations, 0)

		assert.Equal(t, 5, len(downMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, downMigrations, migration)
			assert.Equal(t, migration.Number, downMigrations[idx].Number)
		}
	})

	t.Run("calculate down migration with last down migration", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		var expected []*migrator.Migration
		for i := 4; i > 0; i-- {
			expected = append(expected, migrations[i-1])
		}

		lastDownMigration := migrations[4]

		downMigrations := migrator.CalculateDownMigrationsToApply(lastDownMigration, migrations, 0)

		assert.Equal(t, 4, len(downMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, downMigrations, migration)
			assert.Equal(t, migration.Number, downMigrations[idx].Number)
		}
	})

	t.Run("calculate down migration with last down migration and num to migrate", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		var expected []*migrator.Migration
		for i := 4; i > 2; i-- {
			expected = append(expected, migrations[i-1])
		}

		lastDownMigration := migrations[4]

		downMigrations := migrator.CalculateDownMigrationsToApply(lastDownMigration, migrations, 2)

		assert.Equal(t, 2, len(downMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, downMigrations, migration)
			assert.Equal(t, migration.Number, downMigrations[idx].Number)
		}
	})

	t.Run("calculate down migration with last up migration and num to migrate", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		var expected []*migrator.Migration
		for i := 5; i > 3; i-- {
			expected = append(expected, migrations[i-1])
		}

		lastUpMigration := migrations[4]

		downMigrations := migrator.CalculateDownMigrationsToApply(lastUpMigration, migrations, 2)

		assert.Equal(t, 2, len(downMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, downMigrations, migration)
			assert.Equal(t, migration.Number, downMigrations[idx].Number)
		}
	})

	t.Run("calculate down migration with last up migration and num to migrate greater than available migrations", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		var expected []*migrator.Migration
		for i := 5; i > 0; i-- {
			expected = append(expected, migrations[i-1])
		}

		lastUpMigration := migrations[4]

		downMigrations := migrator.CalculateDownMigrationsToApply(lastUpMigration, migrations, 20)

		assert.Equal(t, 5, len(downMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, downMigrations, migration)
			assert.Equal(t, migration.Number, downMigrations[idx].Number)
		}
	})

	t.Run("calculate down migration with last down migration and num to migrate greater than available migrations", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		var expected []*migrator.Migration
		for i := 4; i > 0; i-- {
			expected = append(expected, migrations[i-1])
		}

		lastDownMigration := migrations[4]

		downMigrations := migrator.CalculateDownMigrationsToApply(lastDownMigration, migrations, 20)

		assert.Equal(t, 4, len(downMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, downMigrations, migration)
			assert.Equal(t, migration.Number, downMigrations[idx].Number)
		}
	})
}

func TestCalculateUpMigration(t *testing.T) {
	t.Run("calculate up migration with last up migration", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		var expected []*migrator.Migration
		for i := 6; i <= 10; i++ {
			expected = append(expected, migrations[i-1])
		}

		lastUpMigration := migrations[4]

		upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, migrations, 0)

		assert.Equal(t, 5, len(upMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, upMigrations, migration)
			assert.Equal(t, migration.Number, upMigrations[idx].Number)
		}
	})

	t.Run("calculate up migration with last down migration", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		var expected []*migrator.Migration
		for i := 5; i <= 10; i++ {
			expected = append(expected, migrations[i-1])
		}

		lastUpMigration := migrations[4]

		upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, migrations, 0)

		assert.Equal(t, 6, len(upMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, upMigrations, migration)
			assert.Equal(t, migration.Number, upMigrations[idx].Number)
		}
	})

	t.Run("calculate up migration with last up migration and num to migrate", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		var expected []*migrator.Migration
		for i := 6; i <= 7; i++ {
			expected = append(expected, migrations[i-1])
		}

		lastUpMigration := migrations[4]

		upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, migrations, 2)

		assert.Equal(t, 2, len(upMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, upMigrations, migration)
			assert.Equal(t, migration.Number, upMigrations[idx].Number)
		}
	})

	t.Run("calculate up migration with last down migration and num to migrate", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		var expected []*migrator.Migration
		for i := 5; i <= 6; i++ {
			expected = append(expected, migrations[i-1])
		}

		lastDownMigration := migrations[4]

		upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, migrations, 2)

		assert.Equal(t, 2, len(upMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, upMigrations, migration)
			assert.Equal(t, migration.Number, upMigrations[idx].Number)
		}
	})

	t.Run("calculate up migration with last up is last entry and num to migrate greater than available migrations", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		lastUpMigration := migrations[9]

		upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, migrations, 20)

		assert.Equal(t, 0, len(upMigrations))
	})

	t.Run("calculate up migration with last down is last entry and num to migrate greater than available migrations", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		lastUpMigration := migrations[9]

		upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, migrations, 0)

		assert.Equal(t, 1, len(upMigrations))

		assert.Equal(t, migrations[9], upMigrations[0])
	})

	t.Run("calculate up migration with last up migration and num to migrate greater than available migrations", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationUp,
			})
		}

		var expected []*migrator.Migration
		for i := 6; i <= 10; i++ {
			expected = append(expected, migrations[i-1])
		}

		lastUpMigration := migrations[4]

		upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, migrations, 20)

		assert.Equal(t, 5, len(upMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, upMigrations, migration)
			assert.Equal(t, migration.Number, upMigrations[idx].Number)
		}
	})

	t.Run("calculate up migration with last down migration and num to migrate greater than available migrations", func(t *testing.T) {
		var migrations []*migrator.Migration
		for i := 1; i <= 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   fmt.Sprintf("test_%d", i),
				Type:   migrator.MigrationDown,
			})
		}

		var expected []*migrator.Migration
		for i := 5; i <= 10; i++ {
			expected = append(expected, migrations[i-1])
		}

		lastDownMigration := migrations[4]

		upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, migrations, 20)

		assert.Equal(t, 6, len(upMigrations))

		for idx, migration := range expected {
			// assert.Contains(t, upMigrations, migration)
			assert.Equal(t, migration.Number, upMigrations[idx].Number)
		}
	})

}
