package models

import (
	"time"

	"github.com/google/uuid"
)

type Theatre struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	City         string    `json:"city"`
	Area         string    `json:"area"`
	Address      string    `json:"address"`
	TotalScreens int       `json:"total_screens"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
