package dbrepo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adharshmk96/stk/pkg/migrator"
)

const (
	mysqlTableExists               = "SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?"
	mysqlSelectMigrationEntries    = "SELECT number, name, migtype, created FROM " + sqlMigrationTable + " ORDER BY id ASC"
	mysqlDropMigrationTable        = "DROP TABLE IF EXISTS " + sqlMigrationTable
	mysqlLastAppliedMigrationEntry = "SELECT number, name, migtype, created FROM " + sqlMigrationTable + " ORDER BY id DESC LIMIT 1"
	mysqlInsertMigrationEntry      = "INSERT INTO " + sqlMigrationTable + " (number, name, migtype) VALUES (?, ?)"
	mysqlMigrationSchema           = `CREATE TABLE IF NOT EXISTS migdb_migration (
		id INT AUTO_INCREMENT PRIMARY KEY,
		number INT NOT NULL,
		name VARCHAR(255) NOT NULL,
		migtype VARCHAR(5) NOT NULL,
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
)

type mysqlRepo struct {
	conn *sql.DB
}

func NewMySqlRepo(db *sql.DB) migrator.DatabaseRepo {
	return &mysqlRepo{
		conn: db,
	}
}

func (mys *mysqlRepo) LoadLastAppliedMigration() (*migrator.Migration, error) {
	err := mys.CreateMigrationTableIfNotExists()
	if err != nil {
		return nil, err
	}

	return mys.GetLastAppliedMigration()
}

func (mys *mysqlRepo) CreateMigrationTableIfNotExists() error {
	_, err := mys.conn.Exec(string(mysqlMigrationSchema))
	if err != nil {
		return err
	}
	return nil
}

func (mys *mysqlRepo) GetLastAppliedMigration() (*migrator.Migration, error) {
	var migrationNumber int
	var migrationName string
	var migtype string
	var migCreated string
	err := mys.conn.QueryRow(mysqlLastAppliedMigrationEntry).Scan(&migrationNumber, &migrationName, &migtype, &migCreated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var created time.Time
	created, _ = time.Parse("2006-01-02 15:04:05", migCreated)
	migType := migrator.MigrationType(migtype)

	return &migrator.Migration{
		Number:  migrationNumber,
		Name:    migrationName,
		Type:    migType,
		Created: created,
	}, nil

}

func (mys *mysqlRepo) ApplyMigration(migration *migrator.Migration) error {
	tx, err := mys.conn.Begin()
	if err != nil {
		return err
	}

	migrationQuery := migration.Query
	_, err = tx.Exec(migrationQuery)
	if err != nil {
		return err
	}

	_, err = tx.Exec(mysqlInsertMigrationEntry, migration.Name, migration.Type, time.Now())
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func migrationTableExists(db *sql.DB, tableName string) (bool, error) {
	var count int
	err := db.QueryRow(mysqlTableExists, tableName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (mys *mysqlRepo) LoadMigrations() ([]*migrator.Migration, error) {
	_, err := migrationTableExists(mys.conn, sqlMigrationTable)
	if err != nil {
		return nil, fmt.Errorf("migration table does not exist")
	}

	rows, err := mys.conn.Query(mysqlSelectMigrationEntries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migEntries []*migrator.Migration
	for rows.Next() {
		var migrationNumber int
		var migrationName string
		var migtype string
		var migCreated string
		err = rows.Scan(
			&migrationNumber,
			&migrationName,
			&migtype,
			&migCreated,
		)
		if err != nil {
			return nil, err
		}
		var created time.Time
		created, err = time.Parse("2006-01-02 15:04:05", migCreated)
		if err != nil {
			fmt.Println(err)
		}
		migType := migrator.MigrationType(migtype)
		migEntries = append(migEntries, &migrator.Migration{
			Number:  migrationNumber,
			Name:    migrationName,
			Type:    migType,
			Created: created,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return migEntries, nil
}

func (mys *mysqlRepo) DeleteMigrationTable() error {
	_, err := mys.conn.Exec(mysqlDropMigrationTable)
	if err != nil {
		return err
	}
	return nil
}
