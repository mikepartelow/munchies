package db_test

import "testing"

// func TestFood(t *testing.T) {
// 	dB := mustInitDb(t)
// 	defer dB.Close()

// 	wantName := "butter"
// 	wantNutrients := db.Nutrients{
// 		{Name: "beta broccolitene"},
// 		{Name: "sodium"},
// 	}

// 	err := db.Food{
// 		Name:      wantName,
// 		Nutrients: wantNutrients,
// 	}.WriteTo(dB)
// 	assert.NoError(t, err)

// 	var nut db.Food
// 	err = nut.ReadFrom(dB)
// 	assert.NoError(t, err)
// 	assert.Equal(t, wantName, nut.Name)
// 	assert.Equal(t, wantNutrients, nut.Nutrients)
// }

// func TestFoods(t *testing.T) {
// 	want := []struct {
// 		Name      string
// 		Nutrients db.Nutrients
// 	}{
// 		{
// 			Name: "swamp gas",
// 			Nutrients: db.Nutrients{
// 				{Name: "beta broccolitene"},
// 				{Name: "sodium"},
// 			},
// 		},
// 		{
// 			Name: "potatoes",
// 			Nutrients: db.Nutrients{
// 				{Name: "shells"},
// 			},
// 		},

// 		{
// 			Name: "molybdenum",
// 			Nutrients: db.Nutrients{
// 				{Name: "iron"},
// 				{Name: "cinnamon"},
// 			},
// 		},
// 	}

// 	dB := mustInitDb(t)
// 	defer dB.Close()

// 	for _, w := range want {
// 		err := db.Food{
// 			Name:      w.Name,
// 			Nutrients: w.Nutrients,
// 		}.WriteTo(dB)
// 		assert.NoError(t, err)
// 	}

// 	var foods db.Foods
// 	err := foods.ReadFrom(dB)
// 	assert.NoError(t, err)
// 	assert.Len(t, foods, 3)

// 	for i := range want {
// 		assert.Equal(t, want[i].Name, foods[i].Name)
// 		assert.Equal(t, want[i].Nutrients, foods[i].Nutrients)
// 	}
// }

func TestUniqueFoodName(t *testing.T) {
	t.Skip("FIXME")
}
