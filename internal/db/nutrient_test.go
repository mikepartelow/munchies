package db_test

import (
	"mp/munchies/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNutrient(t *testing.T) {
	dB := mustInitDb(t)
	defer dB.Close()

	wantName := "butter"

	err := db.Nutrient{
		Name: wantName,
	}.WriteTo(dB)
	assert.NoError(t, err)

	var nut db.Nutrient
	err = nut.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Equal(t, wantName, nut.Name)
}

func TestNutrients(t *testing.T) {
	names := []string{
		"butter",
		"cheese",
		"hydrogen",
	}

	dB := mustInitDb(t)
	defer dB.Close()

	for _, name := range names {
		err := db.Nutrient{
			Name: name,
		}.WriteTo(dB)
		assert.NoError(t, err)
	}

	var nuts db.Nutrients
	err := nuts.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Len(t, nuts, 3)

	for i := range names {
		assert.Equal(t, names[i], nuts[i].Name)
	}
}
