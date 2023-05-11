package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PGDatabase struct {
	Host             string
	Port             string
	User             string
	Password         string
	Database         string
	ConnectionString string
}

// NewPGDatabase creates a new PostgreSQL database connection
func NewPGDatabase(host, port, user, password, database string) *PGDatabase {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	return &PGDatabase{
		Host:             host,
		Port:             port,
		User:             user,
		Password:         password,
		Database:         database,
		ConnectionString: connectionString,
	}
}

func (pg *PGDatabase) GetPGConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), pg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("connection to PostgreSQL database failed: %w", err)
	}
	return conn, nil
}

func (pg *PGDatabase) GetPGPool() (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), pg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("connection to PostgreSQL database failed: %w", err)
	}
	return pool, nil
}
