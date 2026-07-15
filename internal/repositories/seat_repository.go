package repositories

import (
	"gorm.io/gorm"

	"github.com/poorna-9/goshow/internal/models"
)

type SeatRepository struct {
	db *gorm.DB
}


func NewSeatRepository (db *gorm.DB ) *SeatRepository{
	return &SeatRepository{
		db : db,
	}
}

func (r *SeatRepository) Create(seat *models.Seat) error {
	return r.db.Create(seat).Error
}

func (r *SeatRepository) FindByScreen(id string) ([]models.Seat, error) {
	var seats []models.Seat
	err := r.db.Where("screen_id = ?", id).Find(&seats).Error
	return seats, err
}