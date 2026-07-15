package models

import (
	"time"

	"github.com/google/uuid"
)

type SeatType string

const (
	SeatRegular SeatType = "regular"
	SeatPremium SeatType = "premium"
	SeatVIP     SeatType = "vip"
)

type Seat struct {
	ID         uuid.UUID
	ScreenID   uuid.UUID
	RowLabel   string // e.g. "A", "B", "C"
	SeatNumber string // e.g. 12  -> combined with RowLabel gives "A12"
	SeatType   SeatType
	CreatedAt  time.Time
}
