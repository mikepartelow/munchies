package main

import (
	"fmt"
	food "mp/munchies/pkg/food"
	nutrient "mp/munchies/pkg/nutrient"
	"os"
	"strings"
)

const (
	PATH_TO_DATA = "./USDA"
)

func main() {
	// FIXME: CLI
	//
	//  ./munchies 'oats' : display list of foods with "oats" in the description, case insensitive
	//  ./munchies 'Oats, whole grain, rolled, old fashioned' : display details for single food, case insensitive
	//  ./munchies 'Oats, whole grain, rolled, old fashioned' '1/4 cup' : details, scaled for 1/4 cup
	//  ./munchies meal.(yaml/json) : display details for a meal (combination of foods)
	//

	// check out the portion conversions: cat FoodData_Central_sr_legacy_food_json_2021-10-28.json| jq -C ".SRLegacyFoods[] |  select(.description==\"Apples, raw, with skin (Includes foods for USDA's Food Distribution Program)\") | .foodPortions"

	// FIXME: food.MustRead("path/to/all/data") let MustRead figure out which datasets to load
	foods := food.MustRead(PATH_TO_DATA)

	portion := "100g" // FIXME: this appears to be an assumption of the data, but probably isn't

	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [search string]", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		term := os.Args[1]
		foods = foods.Filter(term)
	}

	if len(foods) == 1 {
		food := foods[0]
		fmt.Println(food.Description + ": " + portion)
		fmt.Println(strings.Repeat("-", 80))
		// FIXME: template
		for _, nut := range food.FoodNutrients {
			// FIXME: --all-nutrients
			// FIXME: less awkward overall
			if nutrient.ShouldDisplay(nut.Nutrient.Name) {
				fmt.Printf("  %40s: %.2f%s\n", nut.Nutrient.Name, nut.Amount, nut.Nutrient.UnitName)
			}
		}

	} else {
		for _, food := range foods {
			fmt.Println(food.Description)
		}
		fmt.Printf("Found %d foods.\n", len(foods))
	}
}
