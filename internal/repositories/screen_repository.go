package repositories

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/poorna-9/goshow/internal/models"
)

type ScreenRepository struct {
	DB *gorm.DB
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

func NewScreenRepository(db *gorm.DB) *ScreenRepository {
	return &ScreenRepository{DB: db}
}

func (r *ScreenRepository) Create(screen *models.Screen) (uuid.UUID, error) {
	err := r.DB.Create(screen).Error
	if err != nil {
		return uuid.Nil, err
	}

	return screen.ID, nil
}

func (r *ScreenRepository) FindByTheatre(theatreID string) ([]models.Screen, error) {
	var screens []models.Screen
	err := r.DB.Where("theatre_id = ? AND is_active = ?", theatreID, true).Find(&screens).Error
	return screens, err
}

func (r *ScreenRepository) FindByID(id string) (*models.Screen, error) {
	var screen models.Screen
	err := r.DB.First(&screen, "id = ? AND is_active = ?", id, true).Error
	return &screen, err
}

func (r *ScreenRepository) Deactivate(id string) error {
	return r.DB.Model(&models.Screen{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *ScreenRepository) CreateSeats(screenID uuid.UUID, layout []SeatLayoutRow) error {
	var seats []models.Seat
	rowOffset := 0

	for _, block := range layout {
		for row := 0; row < block.Rows; row++ {
			rowLabel := string(rune('A' + rowOffset + row))
			for seatNum := 1; seatNum <= block.SeatsPerRow; seatNum++ {
				seats = append(seats, models.Seat{
					ID:         uuid.New(),
					ScreenID:   screenID,
					RowLabel:   rowLabel,
					SeatNumber: fmt.Sprintf("%s%d", rowLabel, seatNum),
					SeatType:   block.SeatType,
				})
			}
		}
		rowOffset += block.Rows
	}

	return r.DB.Create(&seats).Error
}
