package db

const (
	FOODS_TABLE = "foods"
)

type Food struct {
	Record
	Name      string
	Nutrients []Nutrient
}

type Foods []Food

func (u Food) WriteTo(db *Database) error {
	// FIXME:
	// nutrient.GetIdMap()
	// db.db.Begin()
	// WriteThing(u)
	// WriteThing(&FoodNutrient{idMap[nutrient]})
	// db.db.End()
	return u.WriteThing(u, FOODS_TABLE, db)
}

func (u *Food) ReadFrom(db *Database) error {
	return u.ReadThing(u, FOODS_TABLE, db)
}

func (u *Foods) ReadFrom(db *Database) error {
	return Record{}.ReadThings(u, FOODS_TABLE, db)
}
