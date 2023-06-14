package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

const (
	// IN_MEMORY, passed to New(), creates an in-memory database
	IN_MEMORY = ":memory:"
)

func New(filename string) (*Database, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't create database %q: %w", filename, err)
	}

	table := "nutrients"
	sql := fmt.Sprintf(`
		create table %s (
			id integer not null primary key,
			name text,
			created_at datetime default current_timestamp,
			updated_at datetime default current_timestamp
		);
	`, table) // FIXME: go template embed

	if _, err = db.Exec(sql); err != nil {
		return nil, fmt.Errorf("couldn't create table %q: %w", table, err)
	}

	return &Database{
		db: db,
	}, nil
}

func (db *Database) Close() {
	db.db.Close()
}
