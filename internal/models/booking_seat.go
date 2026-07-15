package models

import (
	"github.com/google/uuid"
)

type BookingSeat struct {
	ID        uuid.UUID
	BookingID uuid.UUID
	SeatID    uuid.UUID
	Price     float64
}
