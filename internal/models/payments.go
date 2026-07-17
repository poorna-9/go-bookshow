package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending        PaymentStatus = "pending"
	PaymentSuccess        PaymentStatus = "success"
	PaymentFailed         PaymentStatus = "failed"
	PaymentRefundRequired PaymentStatus = "refund_required"
)

type Payment struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	SessionID       uuid.UUID
	ShowID          uuid.UUID
	Amount          float64
	Status          PaymentStatus
	Method          string
	TransactionRef  string
	RazorpayOrderID string
	ExpiresAt       time.Time
	CreatedAt       time.Time
}
