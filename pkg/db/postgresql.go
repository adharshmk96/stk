package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgInstance *pgx.Conn
var poolInstance *pgxpool.Pool
var pgOnce sync.Once
var poolOnce sync.Once

func GetPGConnection(ctx context.Context, host, port, database, user, password string) *pgx.Conn {
	pgOnce.Do(func() {
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, database)
		conn, err := pgx.Connect(ctx, connStr)
		if err != nil {
			panic(err)
		}
		pgInstance = conn
	})
	return pgInstance
}

func GetPGPool(host, port, database, user, password string) *pgxpool.Pool {
	poolOnce.Do(func() {
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, database)
		pool, err := pgxpool.New(context.Background(), connStr)
		if err != nil {
			panic(err)
		}
		poolInstance = pool
	})
	return poolInstance
}

func ResetPGConnection() {
	pgInstance = nil
	pgOnce = sync.Once{}
}
