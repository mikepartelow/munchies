package ingest_test

import (
	"mp/munchies/internal/db"
	"mp/munchies/pkg/ingest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngest(t *testing.T) {
	wantNutrients := []string{
		"Ash",
		"Cryptoxanthin, beta",
		"Energy",
		"Protein",
	}
	wantUnits := []string{
		"g",
		"kcal",
		"µg",
	}
	wantFoods := []string{
		"Hummus, commercial",
		"Pillsbury Golden Layer Buttermilk Biscuits, Artificial Flavor, refrigerated dough",
	}

	i, err := ingest.New("testdata", db.IN_MEMORY)
	assert.NoError(t, err)
	defer i.DB.Close()
	assert.Len(t, wantNutrients, i.Nutrients)
	assert.Len(t, wantUnits, i.Units)
	assert.Len(t, wantFoods, i.Foods)

	assert.ElementsMatch(t, wantNutrients, gotNutrients(t, i.DB))
	assert.ElementsMatch(t, wantUnits, gotUnits(t, i.DB))
	assert.ElementsMatch(t, wantFoods, gotFoods(t, i.DB))
}

func TestIngestRelations(t *testing.T) {
	i, err := ingest.New("testdata", db.IN_MEMORY)
	assert.NoError(t, err)
	defer i.DB.Close()

	var food db.Food
	assert.NoError(t, food.ReadFrom(i.DB))
	assert.NotEmpty(t, food.Nutrients)

	var foods db.Foods
	assert.NoError(t, foods.ReadFrom(i.DB))
	assert.NotEmpty(t, foods)
	for _, food := range foods {
		assert.NotEmpty(t, food.Nutrients)
	}

}

func gotNutrients(t *testing.T, dB *db.Database) []string {
	var gotNutrients db.Nutrients
	assert.NoError(t, gotNutrients.ReadFrom(dB))

	nuts := make([]string, len(gotNutrients))
	for i := range gotNutrients {
		nuts[i] = gotNutrients[i].Name
	}

	return nuts
}

func gotUnits(t *testing.T, dB *db.Database) []string {
	var gotUnits db.Units
	assert.NoError(t, gotUnits.ReadFrom(dB))

	units := make([]string, len(gotUnits))
	for i := range gotUnits {
		units[i] = gotUnits[i].Name
	}

	return units
}

func gotFoods(t *testing.T, dB *db.Database) []string {
	var gotFoods db.Foods
	assert.NoError(t, gotFoods.ReadFrom(dB))

	foods := make([]string, len(gotFoods))
	for i := range gotFoods {
		foods[i] = gotFoods[i].Name
	}

	return foods
}
