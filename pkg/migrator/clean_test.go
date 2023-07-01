package migrator_test

import (
	"errors"
	"testing"

	"github.com/adharshmk96/stk/mocks"
	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClean(t *testing.T) {
	t.Run("when database is empty", func(t *testing.T) {

		t.Run("when file system is empty", func(t *testing.T) {
			fsRepo := mocks.NewFileRepo(t)
			dbRepo := mocks.NewDatabaseRepo(t)

			dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return([]*migrator.Migration{}, nil)
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return([]*migrator.Migration{}, nil)

			config := &migrator.MigratorConfig{
				FSRepo: fsRepo,
				DBRepo: dbRepo,
				DryRun: false,
			}

			migrations, err := migrator.Clean(config)

			assert.NoError(t, err)
			assert.Equal(t, 0, len(migrations))
			fsRepo.AssertExpectations(t)
			dbRepo.AssertExpectations(t)
			fsRepo.AssertNotCalled(t, "DeleteMigrationFile", mock.Anything)
		})

		t.Run("when file system has migrations", func(t *testing.T) {
			var upMigrationsFromFile []*migrator.Migration
			for i := 1; i <= 10; i++ {
				upMigrationsFromFile = append(upMigrationsFromFile, &migrator.Migration{
					Number: i,
					Name:   "migration",
					Type:   migrator.MigrationUp,
				})

			}

			var downMigrationsFromFile []*migrator.Migration
			for i := 1; i <= 10; i++ {
				downMigrationsFromFile = append(downMigrationsFromFile, &migrator.Migration{
					Number: i,
					Name:   "migration",
					Type:   migrator.MigrationDown,
				})

			}

			fsRepo := mocks.NewFileRepo(t)
			dbRepo := mocks.NewDatabaseRepo(t)

			dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(upMigrationsFromFile, nil).Once()
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return(downMigrationsFromFile, nil).Once()
			fsRepo.On("DeleteMigrationFile", mock.AnythingOfType("*migrator.Migration")).Return(nil).Times(20)

			config := &migrator.MigratorConfig{
				FSRepo: fsRepo,
				DBRepo: dbRepo,
				DryRun: false,
			}

			migrations, err := migrator.Clean(config)

			assert.NoError(t, err)

			assert.Equal(t, 20, len(migrations))
			fsRepo.AssertExpectations(t)
			dbRepo.AssertExpectations(t)

		})

		t.Run("when file system has migrations and dry run is true", func(t *testing.T) {
			var upMigrationsFromFile []*migrator.Migration
			for i := 1; i <= 10; i++ {
				upMigrationsFromFile = append(upMigrationsFromFile, &migrator.Migration{
					Number: i,
					Name:   "migration",
					Type:   migrator.MigrationUp,
				})

			}

			var downMigrationsFromFile []*migrator.Migration
			for i := 1; i <= 10; i++ {
				upMigrationsFromFile = append(upMigrationsFromFile, &migrator.Migration{
					Number: i,
					Name:   "migration",
					Type:   migrator.MigrationDown,
				})

			}

			fsRepo := mocks.NewFileRepo(t)
			dbRepo := mocks.NewDatabaseRepo(t)

			dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(upMigrationsFromFile, nil).Once()
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return(downMigrationsFromFile, nil).Once()

			config := &migrator.MigratorConfig{
				FSRepo: fsRepo,
				DBRepo: dbRepo,
				DryRun: true,
			}

			migrations, err := migrator.Clean(config)

			assert.NoError(t, err)

			assert.Equal(t, 20, len(migrations))

			fsRepo.AssertExpectations(t)
			dbRepo.AssertExpectations(t)
			fsRepo.AssertNotCalled(t, "DeleteMigrationFile", mock.Anything)

		})
	})

	t.Run("when database has migrations", func(t *testing.T) {

		t.Run("when file system is empty", func(t *testing.T) {
			fsRepo := mocks.NewFileRepo(t)
			dbRepo := mocks.NewDatabaseRepo(t)

			dbRepo.On("LoadLastAppliedMigration").Return(&migrator.Migration{
				Number: 5,
				Name:   "migration",
				Type:   migrator.MigrationUp,
			}, nil)
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return([]*migrator.Migration{}, nil)
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return([]*migrator.Migration{}, nil)

			config := &migrator.MigratorConfig{
				FSRepo: fsRepo,
				DBRepo: dbRepo,
				DryRun: false,
			}

			migrations, err := migrator.Clean(config)

			assert.NoError(t, err)
			assert.Equal(t, 0, len(migrations))
			fsRepo.AssertExpectations(t)
			dbRepo.AssertExpectations(t)
			fsRepo.AssertNotCalled(t, "DeleteMigrationFile", mock.Anything)
		})

		t.Run("when file system has migrations", func(t *testing.T) {
			var upMigrationsFromFile []*migrator.Migration
			for i := 1; i <= 10; i++ {
				upMigrationsFromFile = append(upMigrationsFromFile, &migrator.Migration{
					Number: i,
					Name:   "migration",
					Type:   migrator.MigrationUp,
				})

			}

			var downMigrationsFromFile []*migrator.Migration
			for i := 1; i <= 10; i++ {
				downMigrationsFromFile = append(downMigrationsFromFile, &migrator.Migration{
					Number: i,
					Name:   "migration",
					Type:   migrator.MigrationDown,
				})

			}

			fsRepo := mocks.NewFileRepo(t)
			dbRepo := mocks.NewDatabaseRepo(t)

			dbRepo.On("LoadLastAppliedMigration").Return(&migrator.Migration{
				Number: 5,
				Name:   "migration",
				Type:   migrator.MigrationUp,
			}, nil)
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(upMigrationsFromFile, nil)
			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return(downMigrationsFromFile, nil)
			fsRepo.On("DeleteMigrationFile", mock.AnythingOfType("*migrator.Migration")).Return(nil).Times(10)

			config := &migrator.MigratorConfig{
				FSRepo: fsRepo,
				DBRepo: dbRepo,
				DryRun: false,
			}

			migrations, err := migrator.Clean(config)

			assert.NoError(t, err)

			assert.Equal(t, 10, len(migrations))

			fsRepo.AssertExpectations(t)
			dbRepo.AssertExpectations(t)

		})
	})

	t.Run("when delete migration file returns error", func(t *testing.T) {
		var upMigrationsFromFile []*migrator.Migration
		for i := 1; i <= 10; i++ {
			upMigrationsFromFile = append(upMigrationsFromFile, &migrator.Migration{
				Number: i,
				Name:   "migration",
				Type:   migrator.MigrationUp,
			})

		}

		var downMigrationsFromFile []*migrator.Migration
		for i := 1; i <= 10; i++ {
			upMigrationsFromFile = append(upMigrationsFromFile, &migrator.Migration{
				Number: i,
				Name:   "migration",
				Type:   migrator.MigrationDown,
			})

		}

		fsRepo := mocks.NewFileRepo(t)
		dbRepo := mocks.NewDatabaseRepo(t)

		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(upMigrationsFromFile, nil).Once()
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return(downMigrationsFromFile, nil).Once()
		fsRepo.On("DeleteMigrationFile", mock.AnythingOfType("*migrator.Migration")).Return(errors.New("error")).Times(1)

		config := &migrator.MigratorConfig{
			FSRepo: fsRepo,
			DBRepo: dbRepo,
			DryRun: false,
		}

		migrations, err := migrator.Clean(config)

		assert.Error(t, err)
		assert.Nil(t, migrations)
		fsRepo.AssertExpectations(t)
		dbRepo.AssertExpectations(t)
	})

	t.Run("when database returns error", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		dbRepo := mocks.NewDatabaseRepo(t)

		dbRepo.On("LoadLastAppliedMigration").Return(nil, errors.New("error"))
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return([]*migrator.Migration{}, nil)
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationDown).Return([]*migrator.Migration{}, nil)

		config := &migrator.MigratorConfig{
			FSRepo: fsRepo,
			DBRepo: dbRepo,
			DryRun: false,
		}

		migrations, err := migrator.Clean(config)

		assert.Error(t, err)
		assert.Nil(t, migrations)
		fsRepo.AssertExpectations(t)
		dbRepo.AssertExpectations(t)
		fsRepo.AssertNotCalled(t, "DeleteMigrationFile", mock.Anything)
	})
}
