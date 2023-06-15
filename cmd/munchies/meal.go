package main

import (
	"fmt"
	"mp/munchies/pkg/food"
	"mp/munchies/pkg/meal"
	"mp/munchies/pkg/nutrient"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func makeMealCommand() *cli.Command {
	return &cli.Command{
		Name:      "meal",
		Aliases:   []string{"m"},
		Usage:     "analyze meal",
		ArgsUsage: "<path to meal.yaml>",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 1 {
				cli.ShowSubcommandHelpAndExit(cCtx, 1)
			}

			return doMeal(cCtx.Args().First())
		},
	}
}

func doMeal(mealPath string) error {
	dB, err := getDb()
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error opening database %q: %s.", getDbPath(), err), 1)
	}

	meal, err := meal.Read(mealPath, dB)
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error reading meal %q: %s.", mealPath, err), 1)
	}

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

	return nil
}
