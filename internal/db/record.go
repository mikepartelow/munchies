package db

import (
	"time"
)

type Record struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
