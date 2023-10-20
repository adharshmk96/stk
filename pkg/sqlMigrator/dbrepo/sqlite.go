package dbrepo

import (
	"database/sql"
	"log"
	"time"

	sqlmigrator "github.com/adharshmk96/stk/pkg/sqlMigrator"
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
	return &sqliteDb{
		conn: conn,
	}
}

func (db *sqliteDb) Exec(query string) error {
	_, err := db.conn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (db *sqliteDb) PushHistory(migration *sqlmigrator.MigrationDBEntry) error {
	_, err := db.conn.Exec(`INSERT INTO `+MIGRATION_TABLE_NAME+` (name, direction) VALUES (?, ?)`, migration.Name, migration.Direction)
	if err != nil {
		return err
	}

	return nil
}

func (db *sqliteDb) LoadHistory() ([]*sqlmigrator.MigrationDBEntry, error) {

	rows, err := db.conn.Query(`SELECT id, name, direction, created FROM ` + MIGRATION_TABLE_NAME + ` ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var migrations []*sqlmigrator.MigrationDBEntry

	for rows.Next() {
		var id int
		var name string
		var direction string
		var created time.Time

		err = rows.Scan(&id, &name, &direction, &created)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, &sqlmigrator.MigrationDBEntry{
			Name:      name,
			Direction: direction,
			Created:   created,
		})
	}

	return []*sqlmigrator.MigrationDBEntry{}, nil
}

func (db *sqliteDb) InitMigrationTable() error {
	// create migration table if not exists
	_, err := db.conn.Exec(`CREATE TABLE IF NOT EXISTS ` + MIGRATION_TABLE_NAME + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		direction VARCHAR(4) NOT NULL,
		created DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
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
