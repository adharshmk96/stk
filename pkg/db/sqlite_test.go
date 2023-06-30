package db_test

import (
	"testing"

	"github.com/adharshmk96/stk/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestGetSqliteConnection(t *testing.T) {

	t.Run("get sqlite connection", func(t *testing.T) {
		db := db.GetSqliteConnection("test.db")
		assert.NotNil(t, db)
	})

	t.Run("get sqlite connection singleton", func(t *testing.T) {
		db1 := db.GetSqliteConnection("test.db")
		db2 := db.GetSqliteConnection("test.db")
		assert.Equal(t, db1, db2)
	})

}
