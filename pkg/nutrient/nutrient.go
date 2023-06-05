package nutrient

type Nutrient struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	UnitName string `json:"unitName"`
}

type Filter struct {
	keepNames map[string]struct{}
}

func NewFilter() Filter {
	f := Filter{}

	f.keepNames = make(map[string]struct{})
	keepers := []string{
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

	for _, keeper := range keepers {
		f.keepNames[keeper] = struct{}{}
	}

	return f
}

func (f *Filter) ShouldDisplay(nut string) bool {
	// FIXME: if we have Energy in kcal, don't display Energy in kJ

	_, ok := f.keepNames[nut]

	return ok
}
