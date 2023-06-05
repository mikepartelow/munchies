package meal

import (
	"mp/munchies/pkg/food"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/go-yaml/yaml"
)

// external interface
type Meal struct {
	Portions Portions
}

type Portions []Portion

type Portion struct {
	Food     food.Food
	Amount   float32
	UnitName string
}

// internal interface
// FIXME: move internal interface to ./internal/
type meal struct {
	Kind  string    `yaml:"kind"`
	Foods mealFoods `yaml:"foods"`
}

type mealFoods []mealFood

type mealFood struct {
	Name     string  `yaml:"name"`
	Amount   float32 `yaml:"amount"`
	UnitName string  `yaml:"unitName"`
}

func MustRead(mealPath string, foods food.Foods) Meal {
	file, err := os.Open(mealPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var yamlMeal meal
	err = yaml.NewDecoder(file).Decode(&yamlMeal)
	if err != nil {
		panic(err)
	}

	if yamlMeal.Kind != "meal/v1" {
		log.Fatal().Msgf("uknown meal kind: %q", yamlMeal.Kind)
	}

	var meal Meal
	for _, f := range yamlMeal.Foods {
		food := foods.Match(f.Name)
		if food == nil {
			log.Fatal().Msgf("couldn't match food %q", f.Name)
			break
		}
		meal.Portions = append(meal.Portions, Portion{
			Food:     *food,
			Amount:   f.Amount,
			UnitName: f.UnitName,
		})
	}

	return meal
}
