package migCommands

import (
	"strconv"

	"github.com/adharshmk96/stk/pkg/db"
	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/adharshmk96/stk/pkg/migrator/dbrepo"
	"github.com/spf13/viper"
)

func getNumberFromArgs(args []string, defaultValue int) int {
	if len(args) == 0 {
		return defaultValue
	}
	num, err := strconv.Atoi(args[0])
	if err != nil {
		return defaultValue
	}
	return num
}

func selectDbRepo(database migrator.Database) migrator.DatabaseRepo {
	switch database {
	case migrator.SQLiteDB:
		{
			viper.SetDefault("migrator.sqlite.filepath", "stkmigration.db")
			filepath := viper.GetString("migrator.sqlite.filepath")
			conn := db.GetSqliteConnection(filepath)
			return dbrepo.NewSQLiteRepo(conn)
		}
	default:
		{
			conn := db.GetSqliteConnection("migration.db")
			return dbrepo.NewSQLiteRepo(conn)
		}

	}
}
