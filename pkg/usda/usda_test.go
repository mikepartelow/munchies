package usda_test

import (
	"mp/munchies/pkg/food"
	"mp/munchies/pkg/usda"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustRead(t *testing.T) {
	want := food.Foods{
		{},
		{},
		{},
	}

	got, err := usda.MustRead("testdata")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
