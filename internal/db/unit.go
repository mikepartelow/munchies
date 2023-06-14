package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Unit struct {
	Record
	Name string
}

type Units []Unit

func (u Unit) WriteTo(db *Database) error {
	sql := "INSERT INTO units(name) VALUES (?)"
	r, err := db.db.Exec(sql, u.Name)
	if err != nil {
		return fmt.Errorf("couldn't insert unit %q: %w", u.Name, err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("couldn't get last inserted id: %w", err)
	}

	log.Debug().Int64("rowid", id).Msgf("inserted %q", u.Name)

	return nil
}

func (u *Unit) ReadFrom(db *Database) error {
	rows, err := db.db.Query("SELECT name FROM units LIMIT 1")
	if err != nil {
		return fmt.Errorf("couldn't read unit: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		// FIXME: test this
		return fmt.Errorf("couldn't fetch unit: %e", err)
	}

	err = rows.Scan(&u.Name)
	if err != nil {
		// FIXME: test this
		return fmt.Errorf("couldn't scan unit: %e", err)
	}

	log.Debug().Msgf("u.Name=%q", u.Name)

	return nil
}

func (u *Units) ReadFrom(db *Database) error {
	rows, err := db.db.Query("SELECT name FROM units")
	if err != nil {
		return fmt.Errorf("couldn't read unit: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			// FIXME: test this
			return fmt.Errorf("couldn't scan unit: %e", err)
		}

		// FIXME: be smarter than this
		*u = append(*u, Unit{Name: name})

		log.Debug().Msgf("u.Name=%q", name)
	}

	return nil

}
