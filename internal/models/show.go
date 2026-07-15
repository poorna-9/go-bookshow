package models

import (
	"time"

	"github.com/google/uuid"
)

type Show struct {
	ID           uuid.UUID
	MovieID      uuid.UUID
	ScreenID     uuid.UUID
	Date         time.Time
	StartTime    time.Time
	EndTime      time.Time
	RegularPrice float64 `json:"regular_price"`
	PremiumPrice float64 `json:"premium_price"`
	VIPPrice     float64 `json:"vip_price"`
	CreatedAt    time.Time
}

type ShowSeat struct {
	ID         uuid.UUID
	SeatID     uuid.UUID
	ShowID     uuid.UUID
	SeatNumber string
	SeatPrice  float64
	Available  bool
	SeatType   string
	IsBlocked  bool
}
