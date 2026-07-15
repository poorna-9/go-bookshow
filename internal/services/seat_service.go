package services

import (
	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/repositories"
)

type SeatService struct {
	repo *repositories.SeatRepository
}

func NewSeatService(repo *repositories.SeatRepository) *SeatService {
	return &SeatService{repo: repo}
}

func (s *SeatService) GetSeatsByScreen(id string) ([]models.Seat, error) {
	return s.repo.FindByScreen(id)
}
