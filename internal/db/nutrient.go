package db

const (
	NUTRIENTS_TABLE = "nutrients"
)

type Nutrient struct {
	Record
	Name string
}

type Nutrients []Nutrient

func (n Nutrient) WriteTo(db *Database) error {
	return n.WriteThing(n, NUTRIENTS_TABLE, db)
}

func (n *Nutrient) ReadFrom(db *Database) error {
	return n.ReadThing(n, NUTRIENTS_TABLE, db)
}

func (n *Nutrients) ReadFrom(db *Database) error {
	return Record{}.ReadThings(n, NUTRIENTS_TABLE, db)
}
