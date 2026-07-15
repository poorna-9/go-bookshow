package models

import (
	"time"

	"github.com/google/uuid"
)

type Screen struct {
	ID         uuid.UUID
	TheatreID  uuid.UUID
	ScreenName string
	TotalSeats int
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
