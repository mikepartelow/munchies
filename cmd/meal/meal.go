package main

import (
	"fmt"
	"mp/munchies/pkg/food"
	"mp/munchies/pkg/meal"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PATH_TO_DATA = "./USDA"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	foods := food.MustRead(PATH_TO_DATA)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path to meal.yaml>", os.Args[0])
		os.Exit(1)
	}

	mealPath := os.Args[1]

	meal := meal.MustRead(mealPath, foods)

	for _, p := range meal.Portions {
		fmt.Printf("%.02f %s %s\n", p.Amount, p.UnitName, p.Food.Description)
	}
}
