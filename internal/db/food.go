package db

const (
	FOODS_TABLE = "foods"
)

type Food struct {
	Record
	Name string
}

type Foods []Food

func (u Food) WriteTo(db *Database) error {
	return u.WriteThing(u, FOODS_TABLE, db)
}

func (u *Food) ReadFrom(db *Database) error {
	return u.ReadThing(u, FOODS_TABLE, db)
}

func (u *Foods) ReadFrom(db *Database) error {
	return Record{}.ReadThings(u, FOODS_TABLE, db)
}
