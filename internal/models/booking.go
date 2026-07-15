package models

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
	BookingExpired   BookingStatus = "expired"
)

type Booking struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ShowID      uuid.UUID
	SessionID   uuid.UUID
	Status      BookingStatus
	TotalAmount float64
	SeatIDs     []uuid.UUID `gorm:"serializer:json"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type BookingSession struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	ShowID       uuid.UUID
	Expired      bool
	Success      bool
	SessionSeats []uuid.UUID
}
