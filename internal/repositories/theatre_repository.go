package repositories

import (
	"github.com/poorna-9/goshow/internal/models"
	"gorm.io/gorm"
)

type TheatreRepository struct {
	db *gorm.DB
}

func NewTheatreRepository(db *gorm.DB) *TheatreRepository {
	return &TheatreRepository{
		db: db,
	}
}

func (r *TheatreRepository) Create(theatre *models.Theatre) {
	r.db.Create(theatre)
}

func (r *TheatreRepository) FindByCity(city string) []models.Theatre {
	var theatres []models.Theatre

	r.db.Where("city = ? AND status = ?", city, "active").Find(&theatres)

	return theatres
}

func (r *TheatreRepository) FindById(id string) []models.Theatre {
	var theatres []models.Theatre
	r.db.Where("id = ?", id).First((&theatres))

	return theatres
}
