package db_test

import (
	"mp/munchies/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit(t *testing.T) {
	dB := mustInitDb(t)
	defer dB.Close()

	wantName := "truckload"

	err := db.Unit{
		Name: wantName,
	}.WriteTo(dB)
	assert.NoError(t, err)

	var unit db.Unit
	err = unit.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Equal(t, wantName, unit.Name)
}

func TestUnits(t *testing.T) {
	names := []string{
		"truckload",
		"planeload",
		"boatload",
		"megaparsec",
		"nanometer",
	}

	dB := mustInitDb(t)
	defer dB.Close()

	for _, name := range names {
		err := db.Unit{
			Name: name,
		}.WriteTo(dB)
		assert.NoError(t, err)
	}

	var units db.Units
	err := units.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Len(t, units, 5)

	for i := range names {
		assert.Equal(t, names[i], units[i].Name)
	}
}
