package nutrient_test

import (
	"mp/munchies/pkg/nutrient"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNutrientFilter(t *testing.T) {
	f := nutrient.NewFilter([]string{
		"pop tarts",
		"butter",
	})

	assert.True(t, f.ShouldDisplay("pop tarts"))
	assert.True(t, f.ShouldDisplay("PoP TaRTs"))
	assert.True(t, f.ShouldDisplay("butter"))
	assert.False(t, f.ShouldDisplay("corn"))
	assert.False(t, f.ShouldDisplay("pop corn"))
}
