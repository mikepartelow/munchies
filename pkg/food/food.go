package food

import (
	"mp/munchies/pkg/nutrient"
	"strings"
)

type Foods []Food

type Food struct {
	ID            uint64         `json:"fdcId"`
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

type ByNutrientName []FoodNutrient

func (n ByNutrientName) Len() int           { return len(n) }
func (n ByNutrientName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByNutrientName) Less(i, j int) bool { return n[i].Nutrient.Name < n[j].Nutrient.Name }

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
