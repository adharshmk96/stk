package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
	_ "github.com/mattn/go-sqlite3"
)

// sqlite implementation
type sqliteDb struct {
	conn *sql.DB
}

func NewSQLiteRepo(filePath string) sqlmigrator.DBRepo {
	conn, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	repo := &sqliteDb{
		conn: conn,
	}

	repo.InitMigrationTable()

	return repo
}

func (db *sqliteDb) Exec(query string) error {
	_, err := db.conn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (db *sqliteDb) PushHistory(migration *sqlmigrator.MigrationDBEntry) error {
	_, err := db.conn.Exec(`INSERT INTO `+MIGRATION_TABLE_NAME+` (number, name, direction) VALUES (?, ?, ?)`, migration.Number, migration.Name, migration.Direction)
	if err != nil {
		return err
	}

	return nil
}

func (db *sqliteDb) LoadHistory() ([]*sqlmigrator.MigrationDBEntry, error) {

	rows, err := db.conn.Query(`SELECT * FROM (
		SELECT id, number, name, direction, created FROM ` + MIGRATION_TABLE_NAME + ` ORDER BY id DESC LIMIT 20
	) ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var migrations []*sqlmigrator.MigrationDBEntry

	for rows.Next() {
		var id int
		var number int
		var name string
		var direction string
		var created time.Time

		err = rows.Scan(&id, &number, &name, &direction, &created)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, &sqlmigrator.MigrationDBEntry{
			Number:    number,
			Name:      name,
			Direction: direction,
			Created:   created,
		})
	}

	return migrations, nil
}

func (db *sqliteDb) InitMigrationTable() error {
	// create migration table if not exists
	_, err := db.conn.Exec(`CREATE TABLE IF NOT EXISTS ` + MIGRATION_TABLE_NAME + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		number INTEGER NOT NULL,
		name VARCHAR(255) NOT NULL,
		direction VARCHAR(4) NOT NULL,
		created DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		fmt.Println("error creating migration table")
		return err
	}

	return nil
}

func (db *sqliteDb) DeleteMigrationTable() error {
	_, err := db.conn.Exec("DROP TABLE IF EXISTS " + MIGRATION_TABLE_NAME)
	if err != nil {
		return err
	}
	return nil
}
