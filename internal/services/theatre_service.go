package services

import (
	"strings"

	"github.com/google/uuid"
	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/repositories"
)

type TheatreService struct {
	repo *repositories.TheatreRepository
}

type TheatreInfoResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Area         string `json:"area"`
	Address      string `json:"address"`
	TotalScreens int    `json:"total_screens"`
	Status       string `json:"status"`
}

func NewTheatreService(repo *repositories.TheatreRepository) *TheatreService {
	return &TheatreService{
		repo: repo,
	}
}

func (s *TheatreService) CreateTheatre(theatre *models.Theatre) *models.Theatre {
	theatre.ID = uuid.New()
	theatre.Name = strings.TrimSpace(theatre.Name)
	theatre.City = strings.ToLower(strings.TrimSpace(theatre.City))
	theatre.Area = strings.TrimSpace(theatre.Area)
	theatre.Address = strings.TrimSpace(theatre.Address)
	theatre.Status = "active"
	s.repo.Create(theatre)
	return theatre
}

func (s *TheatreService) GetTheatresByCity(city string) []models.Theatre {
	city = strings.ToLower(strings.TrimSpace(city))
	theatres := s.repo.FindByCity(city)
	return theatres
}

func (s *TheatreService) Gettheatreinfo(id string) TheatreInfoResponse {
	theatres := s.repo.FindById(id)
	if len(theatres) == 0 {
		return TheatreInfoResponse{}
	}
	theatre := theatres[0]
	return TheatreInfoResponse{
		ID:           theatre.ID.String(),
		Name:         theatre.Name,
		City:         theatre.City,
		Area:         theatre.Area,
		Address:      theatre.Address,
		TotalScreens: theatre.TotalScreens,
		Status:       theatre.Status,
	}
}
