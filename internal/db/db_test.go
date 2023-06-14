package db_test

import (
	"mp/munchies/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	dB, err := db.New(db.IN_MEMORY)
	assert.NoError(t, err)
	dB.Close()
}

func TestNutrient(t *testing.T) {
	dB, err := db.New(db.IN_MEMORY)
	assert.NoError(t, err)
	defer dB.Close()

	err = db.Nutrient{
		Name: "butter",
	}.WriteTo(dB)
	assert.NoError(t, err)

	var nut db.Nutrient
	err = nut.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Equal(t, "butter", nut.Name)
}

func TestNutrients(t *testing.T) {
	names := []string{
		"butter",
		"cheese",
		"hydrogen",
	}

	dB, err := db.New(db.IN_MEMORY)
	assert.NoError(t, err)
	defer dB.Close()

	for _, name := range names {
		err = db.Nutrient{
			Name: name,
		}.WriteTo(dB)
		assert.NoError(t, err)
	}

	var nuts db.Nutrients
	err = nuts.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Len(t, nuts, 3)

	for i, _ := range names {
		assert.Equal(t, names[i], nuts[i].Name)
	}

}

func TestRecord(t *testing.T) {
	// test that ID, CreatedAt, UpdatedAt are populated
	t.Skip("FIXME")
	// assert.Equal(t, "butter", nuts[0].ID)
	// assert.Equal(t, "butter", nuts[0].CreatedAt)
	// assert.Equal(t, "butter", nuts[0].UpdatedAt)
}
