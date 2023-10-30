package pingStorage

import (
	"database/sql"
	"fmt"

	"github.com/adharshmk96/stk-template/singlemod/internals/core/entity"
	"github.com/adharshmk96/stk-template/singlemod/internals/core/serr"
	"github.com/adharshmk96/stk-template/singlemod/server/infra"
)

type sqliteRepo struct {
	conn *sql.DB
}

func NewSqliteRepo(conn *sql.DB) entity.PingStorage {
	return &sqliteRepo{
		conn: conn,
	}
}

// Repository Methods
func (s *sqliteRepo) Ping() error {
	rows, err := s.conn.Query("SELECT 1")
	if err != nil {
		return serr.ErrPingFailed
	}
	defer rows.Close()

	var result int

	if rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return serr.ErrPingFailed
		}
	} else {
		return serr.ErrPingFailed
	}

	logger := infra.GetLogger()
	logger.Info(fmt.Sprintf("connection result: %d", result))
	return nil
}
