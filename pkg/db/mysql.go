package db

import (
	"database/sql"
	"fmt"
	"sync"
)

var (
	mysqlConn *sql.DB
	mysqlOnce sync.Once
)

func GetMysqlConnection(host, port, dbname, user, password string) *sql.DB {
	mysqlOnce.Do(func() {
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
		conn, err := sql.Open("mysql", connStr)
		if err != nil {
			panic(err)
		}
		mysqlConn = conn
	})
	return mysqlConn
}
