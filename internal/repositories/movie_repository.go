package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/poorna-9/goshow/internal/models"
	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

type ShowWithTheatre struct {
	ShowID      uuid.UUID
	StartTime   time.Time
	EndTime     time.Time
	BasePrice   float64
	ScreenName  string
	TheatreID   uuid.UUID
	TheatreName string
	TheatreArea string
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) Create(movie *models.Movie) error {
	return r.db.Create(movie).Error
}

func (r *MovieRepository) FindAll() ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Find(&movies).Error
	return movies, err
}

func (r *MovieRepository) FindById(id string) ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Where("id = ?", id).Find(&movies).Error
	return movies, err
}

func (r *MovieRepository) FindByCity(city string) ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Distinct().
		Joins("JOIN shows on shows.movie_id = movies.id").
		Joins("JOIN screens on screens.id = shows.screen_id").
		Joins("JOIN theatres on theatres.id = screens.theatre_id").
		Where("theatres.city = ?", city).
		Find(&movies).Error
	return movies, err

}
