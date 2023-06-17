package nutrient

import "strings"

type Nutrient struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	UnitName string `json:"unitName"`
}

var (
	NUTRIENT_FILTER_DEFAULTS = []string{
		"Cholesterol",
		"Caffeine",
		"Sugars, total including NLEA",
		"Sodium, Na",
		"Fiber, total dietary",
		"Protein",
		"Energy",
		"Fatty acids, total saturated",
		"Total lipid (fat)",
		"Total dietary fiber (AOAC 2011.25)",
		"Carbohydrate, by difference",
		"Energy (Atwater General Factors)",
		"Energy (Atwater Specific Factors)",
	}
)

type Filter struct {
	keepNames map[string]struct{}
}

func NewFilter(keepers []string) Filter {
	f := Filter{}

	f.keepNames = make(map[string]struct{})

	for _, keeper := range keepers {
		f.keepNames[strings.ToLower(keeper)] = struct{}{}
	}

	return f
}

func (f *Filter) ShouldDisplay(nut string) bool {
	// FIXME: if we have Energy in kcal, don't display Energy in kJ

	_, ok := f.keepNames[strings.ToLower(nut)]

	return ok
}
