package main

import (
	"fmt"
	"mp/munchies/pkg/usda"
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

	foods, err := usda.MustRead(PATH_TO_DATA)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	units := make(map[string]struct{})

	for _, f := range foods {
		for _, fn := range f.FoodNutrients {
			units[fn.Nutrient.UnitName] = struct{}{}
		}
	}

	for unit := range units {
		fmt.Println(unit)
	}
}
