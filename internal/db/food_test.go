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
	wants, dB := makeFoods(t)
	defer dB.Close()

	var foods db.Foods
	err := foods.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Len(t, foods, 3)

	for i := range wants {
		assert.Equal(t, wants[i].Name, foods[i].Name)
		assert.Equal(t, scrubTimestamps(wants[i].Nutrients), scrubTimestamps(foods[i].Nutrients))
	}
}

func TestMatchOne(t *testing.T) {
	wants, dB := makeFoods(t)
	defer dB.Close()

	for _, want := range wants {
		var got db.Foods
		assert.NoError(t, got.Match(want.Name, dB))
		assert.Len(t, got, 1)
		assert.Equal(t, want.Name, got[0].Name)
		assert.Len(t, got[0].Nutrients, len(want.Nutrients))

		got = db.Foods{}
		partial := want.Name[2 : len(want.Name)-2]
		assert.NoError(t, got.Match(partial, dB))
		assert.Len(t, got, 1)
		assert.Equal(t, want.Name, got[0].Name)
		assert.Len(t, got[0].Nutrients, len(want.Nutrients))
	}

	var got db.Foods
	assert.Equal(t, db.NoMatch, got.Match("this is not in the db", dB))
}

func TestMatchSeveral(t *testing.T) {
	wants, dB := makeFoods(t)
	defer dB.Close()

	var got db.Foods
	assert.NoError(t, got.Match("a", dB))
	assert.Len(t, got, 2)

	assert.Equal(t, wants[0].Name, got[0].Name)
	assert.Equal(t, wants[1].Name, got[1].Name)

	assert.Len(t, got[0].Nutrients, 2)
	assert.Len(t, got[1].Nutrients, len(wants[1].Nutrients))
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

type foodFixture struct {
	ID        uint64
	Name      string
	Nutrients db.Nutrients
}

func makeFoods(t *testing.T) ([]foodFixture, *db.Database) {
	t.Helper()

	wants := []foodFixture{
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

	for _, w := range wants {
		err := db.Food{
			Record:    db.Record{ID: w.ID},
			Name:      w.Name,
			Nutrients: w.Nutrients,
		}.WriteTo(dB)
		assert.NoError(t, err)
	}

	return wants, dB
}
