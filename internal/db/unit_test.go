package db_test

import (
	"mp/munchies/internal/db"
	"sort"
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
		"boatload",
		"megaparsec",
		"nanometer",
		"planeload",
		"truckload",
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

	sort.Slice(units, func(i, j int) bool {
		return units[i].Name < units[j].Name
	})

	for i := range names {
		assert.Equal(t, names[i], units[i].Name)
	}
}

func TestUniqueUnitName(t *testing.T) {
	t.Skip("FIXME")
}
