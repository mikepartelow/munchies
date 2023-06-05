package nutrient

type Nutrient struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	UnitName string `json:"unitName"`
}

func ShouldDisplay(nut string) bool {
	// FIXME: use a map
	// FIXME: if we have Energy in kcal, don't display Energy in kJ
	want := []string{
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

	for _, w := range want {
		if nut == w {
			return true
		}
	}

	return false
}
