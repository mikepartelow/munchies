package db_test

import (
	"mp/munchies/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	dB, err := db.New(db.IN_MEMORY)
	assert.NoError(t, err)
	defer dB.Close()

	version, err := dB.Version()
	assert.NoError(t, err)
	assert.Equal(t, "0000", version)
}

func TestRecord(t *testing.T) {
	// test that ID, CreatedAt, UpdatedAt are populated
	t.Skip("FIXME")
	// assert.Equal(t, "butter", nuts[0].ID)
	// assert.Equal(t, "butter", nuts[0].CreatedAt)
	// assert.Equal(t, "butter", nuts[0].UpdatedAt)
}

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

func mustInitDb(t *testing.T) *db.Database {
	t.Helper()

	dB, err := db.New(db.IN_MEMORY)
	assert.NoError(t, err)
	assert.NoError(t, dB.Migrate())

	return dB
}
