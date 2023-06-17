package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type Database struct {
	db *sql.DB
}

const (
	// IN_MEMORY, passed to New(), creates an in-memory database
	IN_MEMORY = ":memory:"
)

var NoMatch = errors.New("no match")

//go:embed sql
var sqlFs embed.FS

func New(filename string) (*Database, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't create/open database %q: %w", filename, err)
	}

	return &Database{db: db}, nil
}

func (db Database) Version() (string, error) {
	rows, err := db.db.Query("SELECT version FROM _meta LIMIT 1")
	if err != nil {
		return "", fmt.Errorf("couldn't read version: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("couldn't fetch version")
	}

	var version string
	err = rows.Scan(&version)
	if err != nil {
		return "", fmt.Errorf("couldn't scan version: %w", err)
	}

	log.Debug().Msgf("version=%q", version)

	return version, nil
}

func (db *Database) Migrate() error {
	var version string
	version, err := db.Version()
	if err != nil {
		version = ""
	}

	entries, err := sqlFs.ReadDir("sql")
	if err != nil {
		return fmt.Errorf("couldn't read embed FS, that's weird: %w", err)
	}

	sort.Slice(
		entries,
		func(i, j int) bool {
			return entries[i].Name() < entries[j].Name()
		},
	)

	for _, entry := range entries {
		parts := strings.Split(entry.Name(), "_")
		migrationVersion := parts[0]

		// FIXME: dry
		if migrationVersion > version {
			if err := applySql(entry.Name(), db.db); err != nil {
				return fmt.Errorf("couldn't apply sql: %w", err)
			}
		}
	}

	return nil
}

func (db *Database) Close() {
	db.db.Close()
}

func applySql(filename string, db *sql.DB) error {
	sql, err := sqlFs.ReadFile(path.Join("sql", filename))
	if err != nil {
		return fmt.Errorf("couldn't read sql %q: %w", filename, err)
	}

	if _, err = db.Exec(string(sql)); err != nil {
		return fmt.Errorf("couldn't run sql %q: %w", filename, err)
	}

	return nil
}
