package usda

import (
	"encoding/json"
	"fmt"
	"mp/munchies/pkg/food"
	"os"
)

var _JSON_FILENAMES = []string{
	"foundationDownload.json",
	"FoodData_Central_sr_legacy_food_json_2021-10-28.json",
}

func MustRead(root string) (food.Foods, error) {
	var foods food.Foods
	for _, filename := range _JSON_FILENAMES {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("couldn't open %q: %w", filename, err)
		}
		defer file.Close()

		var foodsMap map[string]food.Foods

		if err = json.NewDecoder(file).Decode(&foodsMap); err != nil {
			return nil, fmt.Errorf("couldn't decode %q: %w", filename, err)
		}

		if len(foodsMap) > 1 {
			return nil, fmt.Errorf("unexpected file format in %q", filename)
		}

		for _, v := range foodsMap {
			foods = append(foods, v...)
		}
	}

	return foods, nil
}
