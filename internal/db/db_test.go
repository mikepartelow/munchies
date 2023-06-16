package db_test

import (
	"mp/munchies/internal/db"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	dB, err := db.New(db.IN_MEMORY)
	assert.NoError(t, err)
	defer dB.Close()

	_, err = dB.Version()
	assert.Error(t, err)

	assert.NoError(t, dB.Migrate())

	version, err := dB.Version()
	assert.NoError(t, err)
	assert.Equal(t, "0000", version)
}

func TestOld(t *testing.T) {
	filename := path.Join(t.TempDir(), "test.sqlite")

	dB, err := db.New(filename)
	assert.NoError(t, err)
	dB.Close()

	dB, err = db.New(filename)
	assert.NoError(t, err)
	dB.Close()
}

func mustInitDb(t *testing.T) *db.Database {
	t.Helper()

	dB, err := db.New(db.IN_MEMORY)
	// dB, err := db.New("/tmp/test.db")
	assert.NoError(t, err)
	assert.NoError(t, dB.Migrate())

	return dB
}
