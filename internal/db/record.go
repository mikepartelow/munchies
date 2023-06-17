package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Record struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r Record) WriteThing(thing any, table string, db *Database) error {
	var result sql.Result
	var err error

	name := reflect.ValueOf(thing).FieldByName("Name").String()
	if id := reflect.ValueOf(thing).FieldByName("ID").Uint(); id != 0 {
		sql := fmt.Sprintf("INSERT INTO %s (id,name) VALUES (?,?)", table)
		result, err = db.db.Exec(sql, id, name)
	} else {
		sql := fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", table)
		result, err = db.db.Exec(sql, name)
	}

	if err != nil {
		return fmt.Errorf("couldn't insert thing %q: %w", name, err)
	}

	iid, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("couldn't get last inserted id: %w", err)
	}

	log.Debug().Int64("rowid", iid).Msgf("inserted %q", name)

	return nil
}

func (r Record) ReadThing(thing any, table string, db *Database) error {
	sql := fmt.Sprintf("SELECT id, name, created_at, updated_at FROM %s LIMIT 1", table)
	rows, err := db.db.Query(sql)
	if err != nil {
		return fmt.Errorf("couldn't read thing: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		// FIXME: test this
		return fmt.Errorf("couldn't fetch thing: %e", err)
	}

	var id uint64
	var name string
	var createdAt time.Time
	var updatedAt time.Time
	err = rows.Scan(&id, &name, &createdAt, &updatedAt)
	if err != nil {
		// FIXME: test this
		return fmt.Errorf("couldn't scan thing: %e", err)
	}

	reflect.Indirect(reflect.ValueOf(thing)).FieldByName("ID").SetUint(id)
	reflect.Indirect(reflect.ValueOf(thing)).FieldByName("Name").SetString(name)
	reflect.Indirect(reflect.ValueOf(thing)).FieldByName("CreatedAt").Set(reflect.ValueOf(createdAt))
	reflect.Indirect(reflect.ValueOf(thing)).FieldByName("UpdatedAt").Set(reflect.ValueOf(updatedAt))

	log.Debug().Msgf("thing=%#v", thing)

	return nil
}

func (r Record) ReadThings(things interface{}, table string, db *Database) error {
	sql := fmt.Sprintf("SELECT id, name, created_at, updated_at FROM %s", table)
	rows, err := db.db.Query(sql)
	if err != nil {
		return fmt.Errorf("couldn't read thing: %w", err)
	}
	defer rows.Close()

	return r.readThings(things, rows)
}

func (r Record) MatchThing(things interface{}, table string, db *Database, term string) error {
	sql := fmt.Sprintf("SELECT id, name, created_at, updated_at FROM %s WHERE name LIKE ?", table)
	rows, err := db.db.Query(sql, "%"+strings.ToLower(term)+"%")
	if err != nil {
		return fmt.Errorf("couldn't read thing: %w", err)
	}
	defer rows.Close()

	err = r.readThings(things, rows)
	if err == nil && reflect.ValueOf(things).Elem().Len() == 0 {
		return NoMatch
	}

	return err
}

func (r Record) readThings(things interface{}, rows *sql.Rows) error {
	thingsSlice := reflect.ValueOf(things).Elem()
	thingType := thingsSlice.Type().Elem()

	for rows.Next() {
		var id uint64
		var name string
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(&id, &name, &createdAt, &updatedAt)
		if err != nil {
			return fmt.Errorf("couldn't scan thing: %w", err)
		}

		newThing := reflect.New(thingType).Elem()

		newThing.FieldByName("ID").SetUint(id)
		newThing.FieldByName("Name").SetString(name)
		newThing.FieldByName("CreatedAt").Set(reflect.ValueOf(createdAt))
		newThing.FieldByName("UpdatedAt").Set(reflect.ValueOf(updatedAt))

		thingsSlice.Set(reflect.Append(thingsSlice, newThing))
	}

	return nil
}
