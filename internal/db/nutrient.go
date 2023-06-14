package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Nutrient struct {
	Record
	Name string
}

type Nutrients []Nutrient

func (n Nutrient) WriteTo(db *Database) error {
	sql := "INSERT INTO nutrients(name) VALUES (?)"
	r, err := db.db.Exec(sql, n.Name)
	if err != nil {
		return fmt.Errorf("couldn't insert nutrient %q: %w", n.Name, err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("couldn't get last inserted id: %w", err)
	}

	log.Debug().Int64("rowid", id).Msgf("inserted %q", n.Name)

	return nil
}

func (n *Nutrient) ReadFrom(db *Database) error {
	rows, err := db.db.Query("SELECT name FROM nutrients LIMIT 1")
	if err != nil {
		return fmt.Errorf("couldn't read nutrient: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		// FIXME: test this
		return fmt.Errorf("couldn't fetch nutrient: %e", err)
	}

	err = rows.Scan(&n.Name)
	if err != nil {
		// FIXME: test this
		return fmt.Errorf("couldn't scan nutrient: %e", err)
	}

	log.Debug().Msgf("n.Name=%q", n.Name)

	return nil
}

func (n *Nutrients) ReadFrom(db *Database) error {
	rows, err := db.db.Query("SELECT name FROM nutrients")
	if err != nil {
		return fmt.Errorf("couldn't read nutrient: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			// FIXME: test this
			return fmt.Errorf("couldn't scan nutrient: %e", err)
		}

		// FIXME: be smarter than this
		*n = append(*n, Nutrient{Name: name})

		log.Debug().Msgf("n.Name=%q", name)
	}

	return nil

}
