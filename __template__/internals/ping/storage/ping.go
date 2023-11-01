package storage

import (
	"database/sql"
	"fmt"

	"github.com/adharshmk96/stktemplate/internals/ping/domain"
	"github.com/adharshmk96/stktemplate/internals/ping/serr"
	"github.com/adharshmk96/stktemplate/server/infra"
)

type sqliteRepo struct {
	conn *sql.DB
}

func NewSqliteRepo(conn *sql.DB) domain.PingStorage {
	return &sqliteRepo{
		conn: conn,
	}
}

func (s *sqliteRepo) Ping() error {
	logger := infra.GetLogger()
	rows, err := s.conn.Query(SELECT_ONE_TEST)
	if err != nil {
		return serr.ErrPingFailed
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error("connection close failed.")
		}
	}(rows)

	var result int

	if rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return serr.ErrPingFailed
		}
	} else {
		return serr.ErrPingFailed
	}

	logger.Info(fmt.Sprintf("connection result: %d", result))
	return nil
}
