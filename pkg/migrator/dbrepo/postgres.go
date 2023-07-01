package dbrepo

import (
	"context"
	"fmt"
	"log"

	"github.com/adharshmk96/stk/pkg/migrator"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	pgMigrationTableExists      = "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)"
	pgLastAppliedMigrationEntry = "SELECT number, name, migtype, created FROM " + sqlMigrationTable + " ORDER BY id DESC LIMIT 1"
	pgInsertMigrationEntry      = "INSERT INTO " + sqlMigrationTable + " (number, name, migtype) VALUES ($1, $2, $3)"
	pgSelectMigrationEntries    = "SELECT number, name, migtype, created FROM " + sqlMigrationTable + " ORDER BY id ASC"
	pgDropMigrationTable        = "DROP TABLE IF EXISTS " + sqlMigrationTable
	pgMigrationSchema           = `CREATE TABLE IF NOT EXISTS migdb_migration (
		id serial primary key,
		number int not null,
		name varchar(255) not null,
		migtype varchar(5) not null,
		created timestamp not null default now()
	);`
)

type postgresRepo struct {
	conn *pgx.Conn
}

func NewPGRepo(conn *pgx.Conn) migrator.DatabaseRepo {
	return &postgresRepo{
		conn: conn,
	}
}

func (pg *postgresRepo) LoadLastAppliedMigration() (*migrator.Migration, error) {
	err := pg.CreateMigrationTableIfNotExists()
	if err != nil {
		return nil, err
	}

	return pg.GetLastAppliedMigration()
}

func (pg *postgresRepo) CreateMigrationTableIfNotExists() error {
	_, err := pg.conn.Exec(context.Background(), string(pgMigrationSchema))
	if err != nil {
		return err
	}
	return nil
}

func (pg *postgresRepo) GetLastAppliedMigration() (*migrator.Migration, error) {
	var migrationNumber int
	var migrationName string
	var migtype string
	var migCreated pgtype.Timestamp
	err := pg.conn.QueryRow(context.Background(), pgLastAppliedMigrationEntry).Scan(&migrationNumber, &migrationName, &migtype, &migCreated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	migType := migrator.MigrationType(migtype)
	return &migrator.Migration{
		Number:  migrationNumber,
		Name:    migrationName,
		Type:    migType,
		Created: migCreated.Time,
	}, nil

}

func (pg *postgresRepo) ApplyMigration(migration *migrator.Migration) error {

	ctx := context.Background()
	tx, err := pg.conn.Begin(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}

	migrationQuery := migration.Query

	_, err = tx.Exec(context.Background(), pgInsertMigrationEntry, migration.Number, migration.Name, migration.Type)
	if err != nil {
		tx.Rollback(ctx)
		log.Fatal(err)
		return err
	}

	_, err = tx.Exec(ctx, migrationQuery)
	if err != nil {
		tx.Rollback(ctx)
		log.Fatal(err)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		tx.Rollback(ctx)
		log.Fatal(err)
		return err
	}

	return nil
}

func (pg *postgresRepo) LoadMigrations() ([]*migrator.Migration, error) {
	var exists bool
	err := pg.conn.QueryRow(context.Background(), pgMigrationTableExists, sqlMigrationTable).Scan(&exists)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("migration table does not exist")
	}

	rows, err := pg.conn.Query(context.Background(), pgSelectMigrationEntries)
	if err != nil {
		return nil, err
	}

	var history []*migrator.Migration

	for rows.Next() {
		var migrationNumber int
		var migrationName string
		var migtype string
		var migCreated pgtype.Timestamp
		err := rows.Scan(&migrationNumber, &migrationName, &migtype, &migCreated)
		if err != nil {
			return nil, err
		}
		migType := migrator.MigrationType(migtype)
		history = append(history, &migrator.Migration{
			Number:  migrationNumber,
			Name:    migrationName,
			Type:    migType,
			Created: migCreated.Time,
		})
	}

	return history, nil
}

func (pg *postgresRepo) DeleteMigrationTable() error {
	_, err := pg.conn.Exec(context.Background(), pgDropMigrationTable)
	if err != nil {
		return err
	}
	return nil
}
