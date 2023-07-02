package tpl

func StorageSqliteTemplate() []byte {
	return []byte(`package sqlite

import (
	"database/sql"

	"github.com/adharshmk96/stk-project-template/pkg/entity"
)

type sqliteRepo struct {
	conn *sql.DB
}

func NewSqliteRepo(conn *sql.DB) entity.PingStorage {
	return &sqliteRepo{
		conn: conn,
	}
}	
`)
}

func StorageSqlitePingTemplate() []byte {
	return []byte(`package sqlite

func (s *sqliteRepo) Ping() error {
	return s.conn.Ping()
}
`)
}
