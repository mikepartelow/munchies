package main

import (
	"fmt"
	"mp/munchies/pkg/food"
	"mp/munchies/pkg/meal"
	"mp/munchies/pkg/nutrient"
	"mp/munchies/pkg/usda"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PATH_TO_DATA = "./USDA"
)

// - SQLITE DB
// - URFAVE CLI
//   ./munchies import /path/to/usda
//   ./munchies units
//   ./munchies meal /path/to/meal
// - TESTS
// - TUI

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	foods, err := usda.MustRead(PATH_TO_DATA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading USDA JSON files at %q.", PATH_TO_DATA)
		log.Fatal().Err(err).Send()
	}

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path to meal.yaml>", os.Args[0])
		os.Exit(1)
	}

	mealPath := os.Args[1]

	meal := meal.MustRead(mealPath, foods)

	for _, p := range meal.Portions {
		fmt.Printf("%.02f %s %s\n", p.Amount, p.UnitName, p.Food.Description)
	}
	fmt.Println(strings.Repeat("-", 80))

	nf := nutrient.NewFilter(nutrient.NUTRIENT_FILTER_DEFAULTS)

	foodNutrients := meal.FoodNutrients()
	sort.Sort(food.ByNutrientName(foodNutrients))

	for _, nut := range foodNutrients {
		if nf.ShouldDisplay(nut.Nutrient.Name) {
			fmt.Printf("  %40s: %.2f%s\n", nut.Nutrient.Name, nut.Amount, nut.Nutrient.UnitName)
		}
	}

}
