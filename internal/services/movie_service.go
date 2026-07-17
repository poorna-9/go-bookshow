package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/repositories"
)

type MovieService struct {
	repo *repositories.MovieRepository
}

func NewMovieService(repo *repositories.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) CreateMovie(movie *models.Movie) error {
	if movie.Title == "" {
		return errors.New("title is required")
	}
	if movie.DurationMinutes <= 0 {
		return errors.New("movie duration should be greater than zero")
	}
	movie.ID = uuid.New()
	return s.repo.Create(movie)
}

func (s *MovieService) GetAllMovies() ([]models.Movie, error) {
	return s.repo.FindAll()
}

func (s *MovieService) GetMovieById(id string) ([]models.Movie, error) {
	return s.repo.FindById(id)
}

func (s *MovieService) GetMovieByCity(city string) ([]models.Movie, error) {
	return s.repo.FindByCity(city)
}

func (s *MovieService) UpdatePosterURL(id uuid.UUID, url string) error {
	return s.repo.UpdatePosterURL(id, url)
}
