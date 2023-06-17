package ingest

import (
	"fmt"
	"mp/munchies/internal/db"
	"mp/munchies/pkg/food"
	"mp/munchies/pkg/usda"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type Ingest struct {
	// A database full of freshly ingested data. Caller must call Close()
	DB *db.Database

	// A count of injested units
	Units int

	// A count of injested nutrients
	Nutrients int

	// A count of injested foods
	Foods int
}

func New(usdaJsonPath, dbPath string) (*Ingest, error) {
	foods, err := usda.MustRead(usdaJsonPath)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, fmt.Errorf("Error reading USDA JSON files at %q: %w", usdaJsonPath, err)
	}

	dB, err := doRecreateDb(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Error rereating database %q: %w", dbPath, err)
	}
	// caller must call dB.Close()

	i, err := doLoadData(dB, foods)
	if err != nil {
		return nil, fmt.Errorf("Error loading data to %q: %w", dbPath, err)
	}

	return i, nil
}

func doRecreateDb(dbPath string) (*db.Database, error) {
	_ = os.RemoveAll(dbPath)
	dB, err := db.New(dbPath)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, fmt.Errorf("Error creating database %q: %w", dbPath, err)
	}
	if err := dB.Migrate(); err != nil {
		log.Error().Err(err).Send()
		return nil, fmt.Errorf("Error migrating database %q: %w", dbPath, err)
	}

	return dB, nil
}

type unitSet map[string]struct{}
type nutrientSet map[string]struct{}
type foodSet map[string]struct{}

func doLoadData(dB *db.Database, foods food.Foods) (*Ingest, error) {
	unitSet, nutrientSet, foodSet := make(unitSet), make(nutrientSet), make(foodSet)

	for _, food := range foods {
		for _, fnut := range food.FoodNutrients {
			if err := doUnit(fnut, unitSet, dB); err != nil {
				return nil, err
			}
			if err := doNutrient(fnut, nutrientSet, dB); err != nil {
				return nil, err
			}
		}
		if err := doFood(food, foodSet, dB); err != nil {
			return nil, err
		}
	}

	return &Ingest{
		DB:        dB,
		Units:     len(unitSet),
		Nutrients: len(nutrientSet),
		Foods:     len(foodSet),
	}, nil
}

func doUnit(fnut food.FoodNutrient, units unitSet, dB *db.Database) error {
	name := fnut.Nutrient.UnitName
	if _, ok := units[name]; !ok {
		if err := (db.Unit{
			Name: strings.TrimSpace(name),
		}.WriteTo(dB)); err != nil {
			log.Error().Err(err).Send()
			return fmt.Errorf("Error writing unit %q to database: %w", name, err)
		}
		units[name] = struct{}{}
	}

	return nil
}

func doNutrient(fnut food.FoodNutrient, nutrients nutrientSet, dB *db.Database) error {
	name := fnut.Nutrient.Name
	if _, ok := nutrients[name]; !ok {
		if err := (db.Nutrient{
			Record: db.Record{
				ID: fnut.Nutrient.ID,
			},
			Name: strings.TrimSpace(name),
		}.WriteTo(dB)); err != nil {
			log.Error().Err(err).Send()
			return fmt.Errorf("Error writing nutrient %q to database: %w", name, err)
		}
		nutrients[name] = struct{}{}
	}

	return nil
}

func doFood(food food.Food, foods foodSet, dB *db.Database) error {
	name := food.Description
	if _, ok := foods[name]; ok {
		log.Warn().Msgf("found duplicate food %q", name)
		return nil
	}

	nutrients := make(db.Nutrients, len(food.FoodNutrients))
	for i := range food.FoodNutrients {
		nutrients[i] = db.Nutrient{
			Record: db.Record{
				ID: food.FoodNutrients[i].Nutrient.ID,
			},
			Name: food.FoodNutrients[i].Nutrient.Name,
		}
	}

	if err := (db.Food{
		Record: db.Record{
			ID: food.ID,
		},
		Name:      strings.TrimSpace(name),
		Nutrients: nutrients,
	}.WriteTo(dB)); err != nil {
		log.Error().Err(err).Send()
		return fmt.Errorf("Error writing food %#v to database: %w", food, err)
	}
	foods[name] = struct{}{}

	return nil
}
