package db

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
)

const (
	FOODS_TABLE = "foods"
)

type Food struct {
	Record
	Name      string
	Nutrients []Nutrient
}

type Foods []Food

func (f Food) WriteTo(db *Database) error {
	// FIXME: transaction!

	for _, nut := range f.Nutrients {
		if err := writeFoodNutrient(f.ID, &nut, db); err != nil {
			return err
		}
	}

	return f.WriteThing(f, FOODS_TABLE, db)
}

func (f *Food) ReadFrom(db *Database) error {
	if err := f.ReadThing(f, FOODS_TABLE, db); err != nil {
		return err
	}
	if f.ID == 0 {
		return fmt.Errorf("unexpected 0 ID")
	}

	nuts, err := readFoodNutrients(f.ID, db)
	if err != nil {
		return err
	}

	f.Nutrients = nuts

	return nil
}

func (f *Foods) ReadFrom(db *Database) error {
	if err := (Record{}.ReadThings(f, FOODS_TABLE, db)); err != nil {
		return err
	}

	for i := range *f {
		nuts, err := readFoodNutrients((*f)[i].ID, db)
		if err != nil {
			return err
		}

		(*f)[i].Nutrients = nuts
	}

	return nil
}

func writeFoodNutrient(food_id uint64, nutrient *Nutrient, db *Database) error {
	if food_id == 0 || nutrient.ID == 0 {
		return fmt.Errorf("got 0 food_id or nutrient_id: %d/%d", food_id, nutrient.ID)
	}
	var result sql.Result
	var err error

	sql := "INSERT INTO nutrients (id,name) VALUES (?,?)"
	_, err = db.db.Exec(sql, nutrient.ID, nutrient.Name)
	if err != nil {
		return fmt.Errorf("couldn't insert nutrient %d/%q: %w", nutrient.ID, nutrient.Name, err)
	}

	sql = "INSERT INTO foods_nutrients (food_id,nutrient_id) VALUES (?,?)"
	result, err = db.db.Exec(sql, food_id, nutrient.ID)
	if err != nil {
		return fmt.Errorf("couldn't insert foods_nutrients %d/%d: %w", food_id, nutrient.ID, err)
	}

	if _, err := result.LastInsertId(); err != nil {
		return fmt.Errorf("couldn't get last inserted id: %w", err)
	}

	log.Debug().Msgf("%d/%d", food_id, nutrient.ID)

	return nil
}

func readFoodNutrients(food_id uint64, db *Database) (Nutrients, error) {
	sql := `SELECT
 		 nutrients.id,
		 nutrients.name
	 FROM foods
	 JOIN foods_nutrients
  	 	ON foods.id=foods_nutrients.food_id
	 JOIN nutrients
  		ON nutrients.id=foods_nutrients.nutrient_id
	 WHERE
	    foods.id=?`

	rows, err := db.db.Query(sql, food_id)
	if err != nil {
		return nil, fmt.Errorf("couldn't read foods_nutrients: %w", err)
	}
	defer rows.Close()

	var nutrients Nutrients
	for rows.Next() {
		var id uint64
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("couldn't scan foods_nutrients: %w", err)
		}
		nutrients = append(nutrients, Nutrient{Record: Record{ID: id}, Name: name})
	}

	return nutrients, nil
}
