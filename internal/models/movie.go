package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID              uuid.UUID
	Title           string
	Description     string
	Genre           string
	Language        string
	DurationMinutes int
	PosterURL       string
	ReleaseDate     time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
