package db_test

import (
	"mp/munchies/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFood(t *testing.T) {
	dB := mustInitDb(t)
	defer dB.Close()

	wantName := "butter"

	err := db.Food{
		Name: wantName,
	}.WriteTo(dB)
	assert.NoError(t, err)

	var nut db.Food
	err = nut.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Equal(t, wantName, nut.Name)
}

func TestFoods(t *testing.T) {
	names := []string{
		"swamp gas",
		"potatoes",
		"molybdenum",
	}

	dB := mustInitDb(t)
	defer dB.Close()

	for _, name := range names {
		err := db.Food{
			Name: name,
		}.WriteTo(dB)
		assert.NoError(t, err)
	}

	var nuts db.Foods
	err := nuts.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Len(t, nuts, 3)

	for i := range names {
		assert.Equal(t, names[i], nuts[i].Name)
	}
}
