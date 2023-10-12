package pingStorage

import (
	"fmt"

	"github.com/adharshmk96/stk-template/multimod/server/infra"
	"github.com/adharshmk96/stk-template/singlemod/internals/core/serr"
)

// Repository Methods
func (s *sqliteRepo) Ping() error {
	res, err := s.conn.Exec("SELECT 1")
	if err != nil {
		return serr.ErrPingFailed
	}
	num, err := res.RowsAffected()
	if err != nil {
		return serr.ErrPingFailed
	}

	logger := infra.GetLogger()
	logger.Info(fmt.Sprintf("Ping Success: %d", num))
	return nil
}
