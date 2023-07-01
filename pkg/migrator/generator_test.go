package migrator_test

import (
	"fmt"
	"testing"

	"github.com/adharshmk96/stk/mocks"
	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerate(t *testing.T) {
	t.Run("generate is succesful with empty directory", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		var migrations []*migrator.Migration
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)
		fsRepo.On("CreateMigrationFile", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		config := migrator.GeneratorConfig{
			Name:          "create_users_table",
			NumToGenerate: 10,
			DryRun:        false,
			Fill:          false,
			FSRepo:        fsRepo,
		}

		err := migrator.Generate(config)
		assert.NoError(t, err)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationsFromFile", 1)
		fsRepo.AssertNumberOfCalls(t, "CreateMigrationFile", config.NumToGenerate*2)
		fsRepo.AssertNotCalled(t, "WriteMigrationToFile")
	})

	t.Run("generate calls write to file when fill is true", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		var migrations []*migrator.Migration
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)
		fsRepo.On("CreateMigrationFile", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		fsRepo.On("WriteMigrationToFile", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		config := migrator.GeneratorConfig{
			Name:          "create_users_table",
			NumToGenerate: 10,
			DryRun:        false,
			Fill:          true,
			FSRepo:        fsRepo,
		}

		err := migrator.Generate(config)
		assert.NoError(t, err)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationsFromFile", 1)
		fsRepo.AssertNumberOfCalls(t, "CreateMigrationFile", config.NumToGenerate*2)
		fsRepo.AssertNumberOfCalls(t, "WriteMigrationToFile", config.NumToGenerate*2)
	})

	t.Run("generate is succesful with non empty directory", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		var migrations []*migrator.Migration
		for i := 0; i < 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   "name",
				Type:   migrator.MigrationUp,
			})
		}
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)
		fsRepo.On("CreateMigrationFile", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		config := migrator.GeneratorConfig{
			Name:          "create_users_table",
			NumToGenerate: 10,
			DryRun:        false,
			Fill:          false,
			FSRepo:        fsRepo,
		}

		err := migrator.Generate(config)
		assert.NoError(t, err)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationsFromFile", 1)
		fsRepo.AssertNumberOfCalls(t, "CreateMigrationFile", config.NumToGenerate*2)
		fsRepo.AssertNotCalled(t, "WriteMigrationToFile")
	})

	t.Run("generate is succesful with dry run", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		var migrations []*migrator.Migration
		for i := 0; i < 10; i++ {
			migrations = append(migrations, &migrator.Migration{
				Number: i,
				Name:   "name",
				Type:   migrator.MigrationUp,
			})
		}
		fsRepo.On("LoadMigrationsFromFile", migrator.MigrationUp).Return(migrations, nil)
		config := migrator.GeneratorConfig{
			Name:          "create_users_table",
			NumToGenerate: 10,
			DryRun:        true,
			Fill:          true,
			FSRepo:        fsRepo,
		}

		err := migrator.Generate(config)
		assert.NoError(t, err)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationsFromFile", 1)
		fsRepo.AssertNotCalled(t, "CreateMigrationFile")
		fsRepo.AssertNotCalled(t, "WriteMigrationToFile")
	})

	t.Run("generate is unsuccesful", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		fsRepo.On("LoadMigrationsFromFile", mock.AnythingOfType("migrator.MigrationType")).Return(nil, fmt.Errorf("error"))
		config := migrator.GeneratorConfig{
			Name:          "create_users_table",
			NumToGenerate: 10,
			DryRun:        false,
			Fill:          false,
			FSRepo:        fsRepo,
		}

		err := migrator.Generate(config)
		assert.Error(t, err)
		fsRepo.AssertNumberOfCalls(t, "LoadMigrationsFromFile", 1)
		fsRepo.AssertNotCalled(t, "CreateMigrationFile")
		fsRepo.AssertNotCalled(t, "WriteMigrationToFile")
	})
}

func TestGenerateMigrationFiles(t *testing.T) {
	t.Run("generate migration files is succesful", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		fsRepo.On("CreateMigrationFile", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		config := migrator.GeneratorConfig{
			Name: "create_users_table",
		}
		migrations := make([]*migrator.Migration, 0)
		for i := 0; i < 10; i++ {
			var migrationType migrator.MigrationType
			if i%2 == 0 {
				migrationType = migrator.MigrationUp
			} else {
				migrationType = migrator.MigrationDown
			}
			migrations = append(migrations, &migrator.Migration{
				Number: i + 1,
				Name:   config.Name,
				Type:   migrationType,
			})
		}

		err := migrator.GenerateMigrationFiles(fsRepo, migrations)
		assert.NoError(t, err)
		fsRepo.AssertNumberOfCalls(t, "CreateMigrationFile", 10)
	})

	t.Run("generate migration files is unsuccesful", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		fsRepo.On("CreateMigrationFile", mock.AnythingOfType("*migrator.Migration")).Return(fmt.Errorf("error"))
		config := migrator.GeneratorConfig{
			Name: "create_users_table",
		}
		migrations := make([]*migrator.Migration, 0)
		for i := 0; i < 10; i++ {
			var migrationType migrator.MigrationType
			if i%2 == 0 {
				migrationType = migrator.MigrationUp
			} else {
				migrationType = migrator.MigrationDown
			}
			migrations = append(migrations, &migrator.Migration{
				Number: i + 1,
				Name:   config.Name,
				Type:   migrationType,
			})
		}

		err := migrator.GenerateMigrationFiles(fsRepo, migrations)
		assert.Error(t, err)
		fsRepo.AssertNumberOfCalls(t, "CreateMigrationFile", 1)
	})
}

func TestFillMigrationFiles(t *testing.T) {
	t.Run("write migration files is succesful", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		fsRepo.On("WriteMigrationToFile", mock.AnythingOfType("*migrator.Migration")).Return(nil)
		config := migrator.GeneratorConfig{
			Name:          "create_users_table",
			NumToGenerate: 10,
			DryRun:        false,
			Fill:          true,
		}
		migrations := make([]*migrator.Migration, 0)
		for i := 0; i < config.NumToGenerate*2; i++ {
			var migrationType migrator.MigrationType
			if i%2 == 0 {
				migrationType = migrator.MigrationUp
			} else {
				migrationType = migrator.MigrationDown
			}
			migrations = append(migrations, &migrator.Migration{
				Number: i + 1,
				Name:   config.Name,
				Type:   migrationType,
			})
		}

		migrator.FillMigrationFiles(fsRepo, config, migrations)

		fsRepo.AssertNumberOfCalls(t, "WriteMigrationToFile", config.NumToGenerate*2)
	})

	t.Run("generate exists if write migration files is unsuccesful", func(t *testing.T) {
		fsRepo := mocks.NewFileRepo(t)
		fsRepo.On("WriteMigrationToFile", mock.AnythingOfType("*migrator.Migration")).Return(fmt.Errorf("error"))
		config := migrator.GeneratorConfig{
			Name:          "create_users_table",
			NumToGenerate: 10,
			DryRun:        false,
			Fill:          true,
		}
		migrations := make([]*migrator.Migration, 0)
		for i := 0; i < config.NumToGenerate*2; i++ {
			var migrationType migrator.MigrationType
			if i%2 == 0 {
				migrationType = migrator.MigrationUp
			} else {
				migrationType = migrator.MigrationDown
			}
			migrations = append(migrations, &migrator.Migration{
				Number: i + 1,
				Name:   config.Name,
				Type:   migrationType,
			})
		}

		err := migrator.FillMigrationFiles(fsRepo, config, migrations)

		assert.Error(t, err)
		fsRepo.AssertNumberOfCalls(t, "WriteMigrationToFile", 1)
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
