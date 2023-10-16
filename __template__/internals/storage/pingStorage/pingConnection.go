package pingStorage

import (
	"database/sql"

	"github.com/adharshmk96/stk-template/singlemod/internals/core/entity"
)

type sqliteRepo struct {
	conn *sql.DB
}

func NewSqliteRepo(conn *sql.DB) entity.PingStorage {
	return &sqliteRepo{
		conn: conn,
	}
}
