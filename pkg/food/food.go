package food

import (
	"encoding/json"
	"mp/munchies/pkg/nutrient"
	"os"
	"path"
	"strings"
)

type Foods []Food

type Food struct {
	FoodClass     string         `json:"foodClass"`
	Description   string         `json:"description"`
	FoodNutrients []FoodNutrient `json:"foodNutrients"`
	// FIXME: Source: which JSON it came from
}

type FoodNutrient struct {
	Type     string            `json:"type"`
	Nutrient nutrient.Nutrient `json:"nutrient"`
	Max      float32           `json:"max"`
	Min      float32           `json:"min"`
	Median   float32           `json:"median"`
	Amount   float32           `json:"amount"`
}

func (fin Foods) Filter(term string) (fout Foods) {
	term = strings.ToLower(term)

	for _, food := range fin {
		if strings.Contains(strings.ToLower(food.Description), term) {
			fout = append(fout, food)
		}
	}

	return
}

func (fin Foods) Match(term string) *Food {
	for _, food := range fin {
		if strings.EqualFold(food.Description, term) {
			return &food
		}
	}

	return nil
}

func MustRead(dataRoot string) (foods Foods) {
	filenames := []string{
		"foundationDownload.json",
		"FoodData_Central_sr_legacy_food_json_2021-10-28.json",
	}

	for _, filename := range filenames {
		dataPath := path.Join(dataRoot, filename)
		file, err := os.Open(dataPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		var foodsMap map[string]Foods

		err = json.NewDecoder(file).Decode(&foodsMap)
		if err != nil {
			panic(err)
		}

		if len(foodsMap) > 1 {
			panic("whoops")
		}

		for _, v := range foodsMap {
			foods = append(foods, v...)
		}
	}

	return
}
