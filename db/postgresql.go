package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var instance *PGDatabase
var once sync.Once
var mutex = &sync.Mutex{}

type PGDatabase struct {
	conn             *pgx.Conn
	pool             *pgxpool.Pool
	ConnectionString string
}

func GetPGDatabaseInstance(host, port, user, password, database string) *PGDatabase {
	once.Do(func() {
		connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, database)

		instance = &PGDatabase{
			ConnectionString: connectionString,
		}
	})
	return instance
}

func (pg *PGDatabase) GetPGConnection() (*pgx.Conn, error) {
	var err error
	pg.conn, err = connectIfNil(pg.conn, pg.ConnectionString)
	return pg.conn, err
}

func (pg *PGDatabase) GetPGPool() (*pgxpool.Pool, error) {
	var err error
	pg.pool, err = poolIfNil(pg.pool, pg.ConnectionString)
	return pg.pool, err
}

func connectIfNil(conn *pgx.Conn, connString string) (*pgx.Conn, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if conn != nil {
		return conn, nil
	}

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("connection to PostgreSQL database failed: %w", err)
	}
	return conn, nil
}

func poolIfNil(pool *pgxpool.Pool, connString string) (*pgxpool.Pool, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if pool != nil {
		return pool, nil
	}

	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("connection to PostgreSQL database failed: %w", err)
	}
	return pool, nil
}
