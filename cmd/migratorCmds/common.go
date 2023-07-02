package migratorCmds

import (
	"context"
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
	case migrator.PostgresDB:
		{
			viper.SetDefault("migrator.postgres.host", "localhost")
			viper.SetDefault("migrator.postgres.port", "5432")
			viper.SetDefault("migrator.postgres.user", "postgres")
			viper.SetDefault("migrator.postgres.password", "postgres")
			viper.SetDefault("migrator.postgres.dbname", "postgres")
			host := viper.GetString("migrator.postgres.host")
			port := viper.GetString("migrator.postgres.port")
			user := viper.GetString("migrator.postgres.user")
			password := viper.GetString("migrator.postgres.password")
			dbname := viper.GetString("migrator.postgres.dbname")

			ctx := context.Background()

			conn := db.GetPGConnection(ctx, host, port, dbname, user, password)
			return dbrepo.NewPGRepo(conn)
		}
	case migrator.MySQLDB:
		{
			viper.SetDefault("migrator.mysql.host", "localhost")
			viper.SetDefault("migrator.mysql.port", "3306")
			viper.SetDefault("migrator.mysql.user", "mysql")
			viper.SetDefault("migrator.mysql.password", "mysql")
			viper.SetDefault("migrator.mysql.dbname", "mysql")
			host := viper.GetString("migrator.mysql.host")
			port := viper.GetString("migrator.mysql.port")
			user := viper.GetString("migrator.mysql.user")
			password := viper.GetString("migrator.mysql.password")
			dbname := viper.GetString("migrator.mysql.dbname")

			conn := db.GetMysqlConnection(host, port, dbname, user, password)
			return dbrepo.NewMySqlRepo(conn)
		}
	default:
		{
			conn := db.GetSqliteConnection("migration.db")
			return dbrepo.NewSQLiteRepo(conn)
		}

	}
}
