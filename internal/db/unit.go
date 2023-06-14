package db

const (
	UNITS_TABLE = "units"
)

type Unit struct {
	Record
	Name string
}

type Units []Unit

func (u Unit) WriteTo(db *Database) error {
	return u.WriteThing(u, UNITS_TABLE, db)
}

func (u *Unit) ReadFrom(db *Database) error {
	return u.ReadThing(u, UNITS_TABLE, db)
}

func (u *Units) ReadFrom(db *Database) error {
	return Record{}.ReadThings(u, UNITS_TABLE, db)
}
