package db

import (
	"fmt"
	"reflect"
	"time"

	"github.com/rs/zerolog/log"
)

type Record struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r Record) WriteThing(thing any, table string, db *Database) error {
	sql := fmt.Sprintf("INSERT INTO %s(name) VALUES (?)", table)
	name := reflect.ValueOf(thing).FieldByName("Name").String()

	result, err := db.db.Exec(sql, name)
	if err != nil {
		return fmt.Errorf("couldn't insert thing %q: %w", name, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("couldn't get last inserted id: %w", err)
	}

	log.Debug().Int64("rowid", id).Msgf("inserted %q", name)

	return nil
}

func (r Record) ReadThing(thing any, table string, db *Database) error {
	sql := fmt.Sprintf("SELECT name FROM %s LIMIT 1", table)
	rows, err := db.db.Query(sql)
	if err != nil {
		return fmt.Errorf("couldn't read thing: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		// FIXME: test this
		return fmt.Errorf("couldn't fetch thing: %e", err)
	}

	var name string
	err = rows.Scan(&name)
	if err != nil {
		// FIXME: test this
		return fmt.Errorf("couldn't scan thing: %e", err)
	}

	reflect.Indirect(reflect.ValueOf(thing)).FieldByName("Name").SetString(name)
	log.Debug().Msgf("thing.Name=%q", name)

	return nil
}

func (r Record) ReadThings(things interface{}, table string, db *Database) error {
	rows, err := db.db.Query(fmt.Sprintf("SELECT name FROM %s", table))
	if err != nil {
		return fmt.Errorf("couldn't read thing: %w", err)
	}
	defer rows.Close()

	thingsSlice := reflect.ValueOf(things).Elem()
	thingType := thingsSlice.Type().Elem()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return fmt.Errorf("couldn't scan thing: %w", err)
		}

		newThing := reflect.New(thingType).Elem()

		newThing.FieldByName("Name").SetString(name)

		thingsSlice.Set(reflect.Append(thingsSlice, newThing))
	}

	return nil
}
