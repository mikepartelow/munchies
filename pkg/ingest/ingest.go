package ingest

import (
	"fmt"
	"mp/munchies/internal/db"
	"mp/munchies/pkg/food"
	"mp/munchies/pkg/usda"
	"os"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

type Ingest struct {
	// A sorted slice of units
	Units []string

	// A sorted slice of nutrients
	Nutrients []string
}

type unitSet map[string]struct{}
type nutrientSet map[string]struct{}

func New(usdaJsonPath, dbPath string) (*Ingest, error) {
	foods, err := usda.MustRead(usdaJsonPath)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, fmt.Errorf("Error reading USDA JSON files at %q: %w", usdaJsonPath, err)
	}
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
	defer dB.Close()

	unitSet, nutrientSet := make(unitSet), make(nutrientSet)
	for _, food := range foods {
		for _, fnut := range food.FoodNutrients {
			if err = doUnit(fnut, unitSet, dB); err != nil {
				return nil, err
			}
			if err = doNutrient(fnut, nutrientSet, dB); err != nil {
				return nil, err
			}
		}
	}

	units, i := make([]string, len(unitSet)), 0
	for u := range unitSet {
		units[i] = u
		i++
	}
	sort.Slice(units, func(i, j int) bool {
		return units[i] < units[j]
	})

	nutrients, i := make([]string, len(nutrientSet)), 0
	for n := range nutrientSet {
		nutrients[i] = n
		i++
	}
	sort.Slice(nutrients, func(i, j int) bool {
		return nutrients[i] < nutrients[j]
	})

	return &Ingest{
		Units:     units,
		Nutrients: nutrients,
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
			Name: strings.TrimSpace(name),
		}.WriteTo(dB)); err != nil {
			log.Error().Err(err).Send()
			return fmt.Errorf("Error writing unit %q to database: %w", name, err)
		}
		nutrients[name] = struct{}{}
	}

	return nil
}
