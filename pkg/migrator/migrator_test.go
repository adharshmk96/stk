package migrator_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/adharshmk96/stk/mocks"
	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

	// TODO: Refactor tests.
	t.Run("when empty database", func(t *testing.T) {

		dbRepo := mocks.NewDatabaseRepo(t)
		dbRepo.On("LoadLastAppliedMigration").Return(nil, nil)

		t.Run("when empty file repo", func(t *testing.T) {
			fsRepo := mocks.NewFileRepo(t)
			var migrations []*migrator.Migration

			fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)

			mockConfig := &migrator.MigratorConfig{
				FSRepo:       fsRepo,
				DBRepo:       dbRepo,
				NumToMigrate: 5,
				DryRun:       false,
			}

			migrations, err := migrator.MigrateUp(mockConfig)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(migrations))

			fsRepo.AssertExpectations(t)
			dbRepo.AssertNotCalled(t, "ApplyMigration")
		})
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
	var upMigrationList []*migrator.Migration
	for i := 1; i <= 10; i++ {
		upMigrationList = append(upMigrationList, &migrator.Migration{
			Number: i,
			Name:   fmt.Sprintf("test_%d", i),
			Type:   migrator.MigrationUp,
		})
	}

	// don't reuse down list because it's pointers
	reversedDownMigrationList := make([]*migrator.Migration, 0)
	for i := 10; i > 0; i-- {
		reversedDownMigrationList = append(reversedDownMigrationList, &migrator.Migration{
			Number: i,
			Name:   fmt.Sprintf("test_%d", i),
			Type:   migrator.MigrationDown,
		})
	}

	var downMigrationList []*migrator.Migration
	for i := 1; i <= 10; i++ {
		downMigrationList = append(downMigrationList, &migrator.Migration{
			Number: i,
			Name:   fmt.Sprintf("test_%d", i),
			Type:   migrator.MigrationDown,
		})
	}

	// length := len(reversedDownMigrationList)

	t.Run("when last migration is up", func(t *testing.T) {

		t.Run("when last migration is at 6", func(t *testing.T) {

			// last migration is 6
			lastMigration := upMigrationList[5]

			tc := []struct {
				name         string
				numToMigrate int
				expected     []*migrator.Migration
			}{
				{
					name:         "when num to migrate is 2",
					numToMigrate: 2,
					// expected should be 6, 5
					expected: reversedDownMigrationList[4:6],
				},
				{
					name:         "when num to migrate is 100",
					numToMigrate: 100,
					// expected should be 6, 5, 4, 3 ,2, 1
					expected: reversedDownMigrationList[4:],
				},
				{
					name:         "when num to migrate is 0",
					numToMigrate: 0,
					// expected should be 6, 5, 4, 3 ,2, 1
					expected: reversedDownMigrationList[4:],
				},
			}

			for _, vtc := range tc {

				t.Run(vtc.name, func(t *testing.T) {

					downMigrations := migrator.CalculateDownMigrationsToApply(lastMigration, reversedDownMigrationList, vtc.numToMigrate)
					assert.Equal(t, len(vtc.expected), len(downMigrations))

					for i, v := range downMigrations {
						assert.Equal(t, vtc.expected[i].Number, v.Number)
						assert.Equal(t, vtc.expected[i].Name, v.Name)
					}
				})
			}

		})

		t.Run("when last migration is at 10", func(t *testing.T) {

			// last migration is 10
			lastMigration := upMigrationList[9]

			tc := []struct {
				name         string
				numToMigrate int
				expected     []*migrator.Migration
			}{
				{
					name:         "when num to migrate is 2",
					numToMigrate: 2,
					// expected should be 10, 9
					expected: reversedDownMigrationList[:2],
				},
				{
					name:         "when num to migrate is 100",
					numToMigrate: 100,
					// expected should be 10, 9, 8, 7, 6, 5, 4, 3, 2, 1
					expected: reversedDownMigrationList,
				},
				{
					name:         "when num to migrate is 0",
					numToMigrate: 0,
					// expected should be 10, 9, 8, 7, 6, 5, 4, 3, 2, 1
					expected: reversedDownMigrationList,
				},
			}

			for _, vtc := range tc {

				t.Run(vtc.name, func(t *testing.T) {

					downMigrations := migrator.CalculateDownMigrationsToApply(lastMigration, reversedDownMigrationList, vtc.numToMigrate)
					assert.Equal(t, len(vtc.expected), len(downMigrations))

					for i, v := range downMigrations {
						assert.Equal(t, vtc.expected[i].Number, v.Number)
						assert.Equal(t, vtc.expected[i].Name, v.Name)
					}
				})
			}
		})

		t.Run("when last migration is at 0", func(t *testing.T) {

			// last migration is 1
			lastMigration := upMigrationList[0]

			tc := []struct {
				name         string
				numToMigrate int
				expected     []*migrator.Migration
			}{
				{
					name:         "when num to migrate is 2",
					numToMigrate: 2,
					// expected should be 0
					expected: reversedDownMigrationList[9:],
				},
				{
					name:         "when num to migrate is 100",
					numToMigrate: 100,
					// expected should be 0
					expected: reversedDownMigrationList[9:],
				},
				{
					name:         "when num to migrate is 0",
					numToMigrate: 0,
					// expected should be 0
					expected: reversedDownMigrationList[9:],
				},
			}

			for _, vtc := range tc {

				t.Run(vtc.name, func(t *testing.T) {

					downMigrations := migrator.CalculateDownMigrationsToApply(lastMigration, reversedDownMigrationList, vtc.numToMigrate)
					assert.Equal(t, len(vtc.expected), len(downMigrations))

					for i, v := range downMigrations {
						assert.Equal(t, vtc.expected[i].Number, v.Number)
						assert.Equal(t, vtc.expected[i].Name, v.Name)
					}
				})
			}
		})

	})

	t.Run("when last migration is down", func(t *testing.T) {

		t.Run("when last migration is at 6", func(t *testing.T) {

			// last migration is 6
			lastMigration := downMigrationList[5]

			tc := []struct {
				name         string
				numToMigrate int
				expected     []*migrator.Migration
			}{
				{
					name:         "when num to migrate is 2",
					numToMigrate: 2,
					// expected should be 5, 4
					expected: reversedDownMigrationList[5:7],
				},
				{
					name:         "when num to migrate is 100",
					numToMigrate: 100,
					// expected should be 5, 4, 3 ,2, 1
					expected: reversedDownMigrationList[5:],
				},
				{
					name:         "when num to migrate is 0",
					numToMigrate: 0,
					// expected should be 5, 4, 3 ,2, 1
					expected: reversedDownMigrationList[5:],
				},
			}

			for _, vtc := range tc {

				t.Run(vtc.name, func(t *testing.T) {

					downMigrations := migrator.CalculateDownMigrationsToApply(lastMigration, reversedDownMigrationList, vtc.numToMigrate)
					assert.Equal(t, len(vtc.expected), len(downMigrations))

					for i, v := range downMigrations {
						assert.Equal(t, vtc.expected[i].Number, v.Number)
						assert.Equal(t, vtc.expected[i].Name, v.Name)
					}
				})
			}

		})

		t.Run("when last migration is at 0", func(t *testing.T) {

			// last migration is 1
			lastMigration := downMigrationList[0]

			tc := []struct {
				name         string
				numToMigrate int
				expected     []*migrator.Migration
			}{
				{
					name:         "when num to migrate is 2",
					numToMigrate: 2,
					// expected should be 0
					expected: []*migrator.Migration{},
				},
				{
					name:         "when num to migrate is 100",
					numToMigrate: 100,
					// expected should be 0
					expected: []*migrator.Migration{},
				},
				{
					name:         "when num to migrate is 0",
					numToMigrate: 0,
					// expected should be 0
					expected: []*migrator.Migration{},
				},
			}

			for _, vtc := range tc {

				t.Run(vtc.name, func(t *testing.T) {

					downMigrations := migrator.CalculateDownMigrationsToApply(lastMigration, reversedDownMigrationList, vtc.numToMigrate)
					assert.Equal(t, len(vtc.expected), len(downMigrations))

					for i, v := range downMigrations {
						assert.Equal(t, vtc.expected[i].Number, v.Number)
						assert.Equal(t, vtc.expected[i].Name, v.Name)
					}
				})
			}
		})

		t.Run("when last migration is at 10", func(t *testing.T) {

			// last migration is 10
			lastMigration := downMigrationList[9]

			tc := []struct {
				name         string
				numToMigrate int
				expected     []*migrator.Migration
			}{
				{
					name:         "when num to migrate is 2",
					numToMigrate: 2,
					// expected should be 9, 8
					expected: reversedDownMigrationList[1:3],
				},
				{
					name:         "when num to migrate is 100",
					numToMigrate: 100,
					// expected should be 9, 8, 7, 6, 5, 4, 3, 2, 1
					expected: reversedDownMigrationList[1:],
				},
				{
					name:         "when num to migrate is 0",
					numToMigrate: 0,
					// expected should be 9, 8, 7, 6, 5, 4, 3, 2, 1
					expected: reversedDownMigrationList[1:],
				},
			}

			for _, vtc := range tc {

				t.Run(vtc.name, func(t *testing.T) {

					downMigrations := migrator.CalculateDownMigrationsToApply(lastMigration, reversedDownMigrationList, vtc.numToMigrate)
					assert.Equal(t, len(vtc.expected), len(downMigrations))

					for i, v := range downMigrations {
						assert.Equal(t, vtc.expected[i].Number, v.Number)
						assert.Equal(t, vtc.expected[i].Name, v.Name)
					}
				})
			}

		})

	})

	t.Run("when last migration is nil", func(t *testing.T) {

		t.Run("when num to migrate is 2", func(t *testing.T) {

			downMigrations := migrator.CalculateDownMigrationsToApply(nil, reversedDownMigrationList, 2)
			assert.Equal(t, 2, len(downMigrations))

			for i, v := range downMigrations {
				assert.Equal(t, reversedDownMigrationList[i].Number, v.Number)
				assert.Equal(t, reversedDownMigrationList[i].Name, v.Name)
			}

		})

		t.Run("when num to migrate is 100", func(t *testing.T) {

			downMigrations := migrator.CalculateDownMigrationsToApply(nil, reversedDownMigrationList, 100)
			assert.Equal(t, 10, len(downMigrations))

			for i, v := range downMigrations {
				assert.Equal(t, reversedDownMigrationList[i].Number, v.Number)
				assert.Equal(t, reversedDownMigrationList[i].Name, v.Name)
			}

		})

		t.Run("when num to migrate is 0", func(t *testing.T) {

			downMigrations := migrator.CalculateDownMigrationsToApply(nil, reversedDownMigrationList, 0)
			assert.Equal(t, 0, len(downMigrations))

		})

	})

	t.Run("when migration list is empty", func(t *testing.T) {

		tc := []struct {
			name         string
			numToMigrate int
			expected     []*migrator.Migration
		}{
			{
				name:         "when num to migrate is 2",
				numToMigrate: 2,
				// expected should be 0
				expected: []*migrator.Migration{},
			},
			{
				name:         "when num to migrate is 100",
				numToMigrate: 100,
				// expected should be 0
				expected: []*migrator.Migration{},
			},
			{
				name:         "when num to migrate is 0",
				numToMigrate: 0,
				// expected should be 0
				expected: []*migrator.Migration{},
			},
		}

		for _, vtc := range tc {

			t.Run(vtc.name, func(t *testing.T) {

				downMigrations := migrator.CalculateDownMigrationsToApply(downMigrationList[5], []*migrator.Migration{}, vtc.numToMigrate)
				assert.Equal(t, len(vtc.expected), len(downMigrations))

			})
		}

		for _, vtc := range tc {

			t.Run(vtc.name, func(t *testing.T) {

				downMigrations := migrator.CalculateDownMigrationsToApply(upMigrationList[2], []*migrator.Migration{}, vtc.numToMigrate)
				assert.Equal(t, len(vtc.expected), len(downMigrations))

			})
		}

	})

}

func TestCalculateUpMigration(t *testing.T) {

	var upMigrationList []*migrator.Migration
	for i := 1; i <= 10; i++ {
		upMigrationList = append(upMigrationList, &migrator.Migration{
			Number: i,
			Name:   fmt.Sprintf("test_%d", i),
			Type:   migrator.MigrationUp,
		})
	}

	var downMigrationList []*migrator.Migration
	for i := 1; i <= 10; i++ {
		downMigrationList = append(downMigrationList, &migrator.Migration{
			Number: i,
			Name:   fmt.Sprintf("test_%d", i),
			Type:   migrator.MigrationDown,
		})
	}

	t.Run("when last migration is up", func(t *testing.T) {

		t.Run("when last migration is at middle", func(t *testing.T) {
			// migration 5
			lastUpMigration := upMigrationList[4]

			t.Run("num to migrate is 2", func(t *testing.T) {
				// migration 6 to 8
				expected := upMigrationList[5:7]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, upMigrationList, 2)

				assert.Equal(t, 2, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}

			})

			t.Run("num to migrate is 100", func(t *testing.T) {
				// migration 6 to 10
				expected := upMigrationList[5:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, upMigrationList, 100)

				assert.Equal(t, 5, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 0", func(t *testing.T) {
				// migration 6 to 10
				expected := upMigrationList[5:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, upMigrationList, 0)

				assert.Equal(t, 5, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

		})

		t.Run("when last migration is at end", func(t *testing.T) {

			// migration 10
			lastUpMigration := upMigrationList[9]

			t.Run("num to migrate is 2", func(t *testing.T) {
				upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, upMigrationList, 2)
				assert.Equal(t, 0, len(upMigrations))
			})

			t.Run("num to migrate is 100", func(t *testing.T) {
				upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, upMigrationList, 100)
				assert.Equal(t, 0, len(upMigrations))
			})

			t.Run("num to migrate is 0", func(t *testing.T) {
				upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, upMigrationList, 0)
				assert.Equal(t, 0, len(upMigrations))
			})
		})

		t.Run("when last migration is at start", func(t *testing.T) {

			// migration 1
			lastUpMigration := upMigrationList[0]

			t.Run("num to migrate is 2", func(t *testing.T) {

				upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, upMigrationList, 2)

				// migration 2 to 3
				expected := upMigrationList[1:3]
				assert.Equal(t, 2, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 100", func(t *testing.T) {
				// migration 2 to 10
				expected := upMigrationList[1:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, upMigrationList, 100)

				assert.Equal(t, 9, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 0", func(t *testing.T) {
				// migration 6 to 10
				expected := upMigrationList[1:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastUpMigration, upMigrationList, 0)

				assert.Equal(t, 9, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})
		})

	})

	t.Run("when last migration is down", func(t *testing.T) {

		t.Run("when last migration is at middle", func(t *testing.T) {
			// migration 5
			lastDownMigration := downMigrationList[4]

			t.Run("num to migrate is 2", func(t *testing.T) {
				// migration 5 to 7
				expected := upMigrationList[4:6]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, upMigrationList, 2)

				assert.Equal(t, 2, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 100", func(t *testing.T) {
				// migration 5 to 10
				expected := upMigrationList[4:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, upMigrationList, 100)

				assert.Equal(t, 6, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 0", func(t *testing.T) {
				// migration 5 to 10
				expected := upMigrationList[4:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, upMigrationList, 0)

				assert.Equal(t, 6, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

		})

		t.Run("when last migration is at end", func(t *testing.T) {

			// migration 10
			lastDownMigration := downMigrationList[9]

			t.Run("num to migrate is 2", func(t *testing.T) {
				// migration 10
				expected := upMigrationList[9:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, upMigrationList, 2)

				assert.Equal(t, 1, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 100", func(t *testing.T) {
				// migration 10
				expected := upMigrationList[9:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, upMigrationList, 100)

				assert.Equal(t, 1, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 0", func(t *testing.T) {
				// migration 10
				expected := upMigrationList[9:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, upMigrationList, 0)

				assert.Equal(t, 1, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})
		})

		t.Run("when last migration is at start", func(t *testing.T) {

			// migration 1
			lastDownMigration := downMigrationList[0]

			t.Run("num to migrate is 2", func(t *testing.T) {

				upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, upMigrationList, 2)

				// migration 2 to 3
				expected := upMigrationList[:2]
				assert.Equal(t, 2, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 100", func(t *testing.T) {
				// migration 2 to 10
				expected := upMigrationList[:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, upMigrationList, 100)

				assert.Equal(t, 10, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 0", func(t *testing.T) {
				// migration 2 to 10
				expected := upMigrationList[:]

				upMigrations := migrator.CalculateUpMigrationsToApply(lastDownMigration, upMigrationList, 0)

				assert.Equal(t, 10, len(upMigrations))

				assert.Equal(t, expected, upMigrations)

				for i, v := range upMigrations {
					assert.Equal(t, expected[i].Name, v.Name)
					assert.Equal(t, expected[i].Number, v.Number)
				}
			})

		})

	})

	t.Run("when last migration is nil", func(t *testing.T) {

		emptyMigrationList := []*migrator.Migration{}

		t.Run("when migrations are not empty", func(t *testing.T) {

			t.Run("num to migrate is 2", func(t *testing.T) {
				upMigrations := migrator.CalculateUpMigrationsToApply(nil, upMigrationList, 2)
				assert.Equal(t, 2, len(upMigrations))
				assert.Equal(t, upMigrationList[:2], upMigrations)
			})

			t.Run("num to migrate is 100", func(t *testing.T) {
				upMigrations := migrator.CalculateUpMigrationsToApply(nil, upMigrationList, 100)
				assert.Equal(t, 10, len(upMigrations))
				for i, v := range upMigrations {
					assert.Equal(t, upMigrationList[i].Name, v.Name)
					assert.Equal(t, upMigrationList[i].Number, v.Number)
				}
			})

			t.Run("num to migrate is 0", func(t *testing.T) {
				upMigrations := migrator.CalculateUpMigrationsToApply(nil, upMigrationList, 0)
				assert.Equal(t, 0, len(upMigrations))
				for i, v := range upMigrations {
					assert.Equal(t, upMigrationList[i].Name, v.Name)
					assert.Equal(t, upMigrationList[i].Number, v.Number)
				}
			})

		})

		t.Run("when migrations are empty", func(t *testing.T) {

			t.Run("num to migrate is 2", func(t *testing.T) {
				upMigrations := migrator.CalculateUpMigrationsToApply(nil, emptyMigrationList, 2)
				assert.Equal(t, 0, len(upMigrations))
			})

			t.Run("num to migrate is 100", func(t *testing.T) {
				upMigrations := migrator.CalculateUpMigrationsToApply(nil, emptyMigrationList, 100)
				assert.Equal(t, 0, len(upMigrations))
			})

			t.Run("num to migrate is 0", func(t *testing.T) {
				upMigrations := migrator.CalculateUpMigrationsToApply(nil, emptyMigrationList, 0)
				assert.Equal(t, 0, len(upMigrations))
			})

		})

	})
}
