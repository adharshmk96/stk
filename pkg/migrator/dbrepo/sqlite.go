package dbrepo

import (
	"database/sql"
	"log"
	"time"

	"github.com/adharshmk96/stk/pkg/migrator"
)

const (
	sqliteMigrationTableExists      = "SELECT 1 FROM sqlite_master WHERE type='table' AND name=?;"
	sqliteSelectMigrationEntries    = "SELECT number, name, migtype, created FROM " + sqlMigrationTable + " ORDER BY id ASC"
	sqliteDropMigrationTable        = "DROP TABLE IF EXISTS " + sqlMigrationTable
	sqliteLastAppliedMigrationEntry = "SELECT number, name, migtype, created FROM " + sqlMigrationTable + " ORDER BY id DESC LIMIT 1"
	sqliteInsertMigrationEntry      = "INSERT INTO " + sqlMigrationTable + " (number, name, migtype) VALUES ($1, $2, $3)"
	sqliteMigrationSchema           = `CREATE TABLE IF NOT EXISTS migdb_migration (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		number INTEGER NOT NULL,
		name VARCHAR(255) NOT NULL,
		migtype VARCHAR(5) NOT NULL,
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
)

type sqliteRepo struct {
	conn *sql.DB
}

func NewSQLiteRepo(conn *sql.DB) migrator.DatabaseRepo {
	return &sqliteRepo{
		conn: conn,
	}
}

func (s *sqliteRepo) LoadLastAppliedMigration() (*migrator.Migration, error) {
	err := s.CreateMigrationTableIfNotExists()
	if err != nil {
		return nil, err
	}

	return s.GetLastAppliedMigration()
}

func (s *sqliteRepo) LoadMigrations() ([]*migrator.Migration, error) {
	var exists bool
	err := s.conn.QueryRow(sqliteMigrationTableExists, sqlMigrationTable).Scan(&exists)

	if err == sql.ErrNoRows {
		return nil, migrator.ErrMigrationTableDoesNotExist
	} else if err != nil {
		return nil, err
	}

	rows, err := s.conn.Query(sqliteSelectMigrationEntries)
	if err != nil {
		return nil, err
	}

	var history []*migrator.Migration

	for rows.Next() {
		var migrationNumber int
		var migrationName string
		var migtype string
		var migCreated time.Time
		err := rows.Scan(
			&migrationNumber,
			&migrationName,
			&migtype,
			&migCreated,
		)
		if err != nil {
			return nil, err
		}
		migType := migrator.MigrationType(migtype)
		history = append(history, &migrator.Migration{
			Number:  migrationNumber,
			Name:    migrationName,
			Type:    migType,
			Created: migCreated,
		})
	}

	return history, nil
}

func (s *sqliteRepo) CreateMigrationTableIfNotExists() error {
	_, err := s.conn.Exec(string(sqliteMigrationSchema))
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteRepo) GetLastAppliedMigration() (*migrator.Migration, error) {
	var migrationNumber int
	var migrationName string
	var migtype string
	var migCreated time.Time
	err := s.conn.QueryRow(sqliteLastAppliedMigrationEntry).Scan(&migrationNumber, &migrationName, &migtype, &migCreated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	migType := migrator.MigrationType(migtype)
	return &migrator.Migration{
		Number:  migrationNumber,
		Name:    migrationName,
		Type:    migType,
		Created: migCreated,
	}, nil
}

func (s *sqliteRepo) ApplyMigration(mig *migrator.Migration) error {

	tx, err := s.conn.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}

	migrationQuery := mig.Query

	_, err = tx.Exec(sqliteInsertMigrationEntry, mig.Number, mig.Name, mig.Type)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	_, err = tx.Exec(migrationQuery)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	return nil
}

func (s *sqliteRepo) DeleteMigrationTable() error {
	_, err := s.conn.Exec(sqliteDropMigrationTable)
	if err != nil {
		return err
	}
	return nil
}
