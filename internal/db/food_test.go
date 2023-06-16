package db_test

import (
	"mp/munchies/internal/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFood(t *testing.T) {
	dB := mustInitDb(t)
	defer dB.Close()

	wantName := "butter"
	wantNutrients := db.Nutrients{
		{Record: db.Record{ID: 1}, Name: "beta broccolitene"},
		{Record: db.Record{ID: 2}, Name: "sodium"},
	}

	err := db.Food{
		Record:    db.Record{ID: 1},
		Name:      wantName,
		Nutrients: wantNutrients,
	}.WriteTo(dB)
	assert.NoError(t, err)

	var nut db.Food
	err = nut.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Equal(t, wantName, nut.Name)
	assert.Equal(t,
		scrubTimestamps(wantNutrients),
		scrubTimestamps(nut.Nutrients))
}

func TestFoods(t *testing.T) {
	want := []struct {
		ID        uint64
		Name      string
		Nutrients db.Nutrients
	}{
		{
			ID:   1,
			Name: "swamp gas",
			Nutrients: db.Nutrients{
				{Record: db.Record{ID: 1}, Name: "beta broccolitene"},
				{Record: db.Record{ID: 2}, Name: "sodium"},
			},
		},
		{
			ID:   2,
			Name: "potatoes",
			Nutrients: db.Nutrients{
				{Record: db.Record{ID: 3}, Name: "shells"},
			},
		},
		{
			ID:   3,
			Name: "molybdenum",
			Nutrients: db.Nutrients{
				{Record: db.Record{ID: 4}, Name: "iron"},
				{Record: db.Record{ID: 5}, Name: "cinnamon"},
			},
		},
	}

	dB := mustInitDb(t)
	defer dB.Close()

	for _, w := range want {
		err := db.Food{
			Record:    db.Record{ID: w.ID},
			Name:      w.Name,
			Nutrients: w.Nutrients,
		}.WriteTo(dB)
		assert.NoError(t, err)
	}

	var foods db.Foods
	err := foods.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Len(t, foods, 3)

	for i := range want {
		assert.Equal(t, want[i].Name, foods[i].Name)
		assert.Equal(t, scrubTimestamps(want[i].Nutrients), scrubTimestamps(foods[i].Nutrients))
	}
}

func TestUniqueFoodName(t *testing.T) {
	t.Skip("FIXME")
}

func scrubTimestamps(nuts db.Nutrients) db.Nutrients {
	for i := range nuts {
		nuts[i].CreatedAt = time.Time{}
		nuts[i].UpdatedAt = time.Time{}
	}

	return nuts
}
