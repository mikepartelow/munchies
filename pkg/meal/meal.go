package meal

import (
	"mp/munchies/pkg/food"
	"os"

	u "github.com/bcicen/go-units"
	"github.com/rs/zerolog/log"

	"github.com/go-yaml/yaml"
)

// external interface
type Meal struct {
	Portions Portions

	foodNutrients []food.FoodNutrient
}

func (m *Meal) FoodNutrients() []food.FoodNutrient {
	panic("insufficient")

	// some essential nutrients (like Energy) are recorded differently for different foods
	// e.g. Oats, whole grain, rolled, old fashioned has no Energy, but has Atwater Specific/General values that must be used instead
	//  so:
	//  - define the nutrients that matter (Fiber, Energy, Fat, etc.)
	//  - compute them if not supplied directly
	//    - kcal instead of kJ
	//    - use Atwater General/Specific for Energy as needed (like in Oats, whole grain, rolled, old fashioned)
	//  - filter out the rest (optionally)
	if m.foodNutrients == nil {
		fnset := make(map[int]food.FoodNutrient)

		for _, p := range m.Portions {
			for _, pffn := range p.Food.FoodNutrients {
				id := pffn.Nutrient.Id
				pffn.Amount *= scale(convert(p))
				if fn, ok := fnset[id]; ok {
					fn.Amount += pffn.Amount
					pffn = fn
				}
				fnset[id] = pffn
			}
		}

		for _, fn := range fnset {
			m.foodNutrients = append(m.foodNutrients, fn)
		}
	}

	return m.foodNutrients
}

// convert converts Portion p into grams
func convert(p Portion) Portion {
	switch p.UnitName {
	case "g":
		break
	case "ounce", "oz":
		p.UnitName = "g"
		p.Amount = float32(u.NewValue(float64(p.Amount), u.Ounce).MustConvert(u.Gram).Float())
	default:
		panic("uknown unit: " + p.UnitName)
	}

	return p
}

// scale returns a scaling factor for a nutrient Amount given that all nutrient values are for 100g of Food.
func scale(p Portion) float32 {
	if p.UnitName != "g" {
		panic("unknown unit: " + p.UnitName)
	}
	return p.Amount / 100.0
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
