package db_test

import (
	"mp/munchies/internal/db"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
	dB := mustInitDb(t)
	defer dB.Close()

	wantName := "truckload"

	err := db.Unit{
		Record: db.Record{
			ID: 2,
		},
		Name: wantName,
	}.WriteTo(dB)
	assert.NoError(t, err)

	var unit db.Unit
	err = unit.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), unit.ID)
	assert.Equal(t, wantName, unit.Name)
	assert.True(t, prettyClose(unit.CreatedAt, time.Now().UTC()))
	assert.True(t, prettyClose(unit.UpdatedAt, time.Now().UTC()))

	var units db.Units
	err = units.ReadFrom(dB)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), units[0].ID)
	assert.Equal(t, wantName, units[0].Name)
	assert.True(t, prettyClose(units[0].CreatedAt, time.Now().UTC()))
	assert.True(t, prettyClose(units[0].UpdatedAt, time.Now().UTC()))
}

func prettyClose(t0, t1 time.Time) bool {
	y0, m0, d0 := t0.Date()
	y1, m1, d1 := t1.Date()

	h0, h1 := t0.Hour(), t1.Hour()

	log.Debug().Msgf("0: %d/%d/%d/%d", y0, m0, d0, h0)
	log.Debug().Msgf("1: %d/%d/%d/%d", y1, m1, d1, h1)

	return y0 == y1 && m0 == m1 && d0 == d1 && h0 == h1
}
