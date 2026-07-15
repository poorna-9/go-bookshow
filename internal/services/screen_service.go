package services

import (
	"github.com/google/uuid"

	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/repositories"
)

type ScreenService struct {
	Repo *repositories.ScreenRepository
}

type SeatLayoutRow struct {
	SeatType    models.SeatType `json:"seat_type"`
	Rows        int             `json:"rows"`
	SeatsPerRow int             `json:"seats_per_row"`
}

type CreateScreenRequest struct {
	TheatreID  uuid.UUID       `json:"theatre_id"`
	ScreenName string          `json:"screen_name"`
	Layout     []SeatLayoutRow `json:"layout"`
}

func NewScreenService(repo *repositories.ScreenRepository) *ScreenService {
	return &ScreenService{Repo: repo}
}

func (s *ScreenService) CreateScreen(req *CreateScreenRequest) error {
	totalSeats := 0
	repoLayout := make([]repositories.SeatLayoutRow, len(req.Layout))
	for i, block := range req.Layout {
		totalSeats += block.Rows * block.SeatsPerRow
		repoLayout[i] = repositories.SeatLayoutRow{
			SeatType:    block.SeatType,
			Rows:        block.Rows,
			SeatsPerRow: block.SeatsPerRow,
		}
	}

	screen := &models.Screen{
		TheatreID:  req.TheatreID,
		ScreenName: req.ScreenName,
		TotalSeats: totalSeats,
		IsActive:   true,
	}

	screenID, err := s.Repo.Create(screen)
	if err != nil {
		return err
	}

	err = s.Repo.CreateSeats(screenID, repoLayout)
	return err
}

func (s *ScreenService) GetScreensByTheatre(theatreID string) ([]models.Screen, error) {
	return s.Repo.FindByTheatre(theatreID)
}

func (s *ScreenService) GetScreenByID(id string) (*models.Screen, error) {
	return s.Repo.FindByID(id)
}

func (s *ScreenService) DeactivateScreen(id string) error {
	return s.Repo.Deactivate(id)
}
