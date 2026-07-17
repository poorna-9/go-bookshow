package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/poorna-9/goshow/internal/models"
	"gorm.io/gorm"
)

type ShowRepository struct {
	db *gorm.DB
}

type ShowTheatre struct {
	ShowID      uuid.UUID
	StartTime   time.Time
	EndTime     time.Time
	BasePrice   float64
	ScreenName  string
	Date        string
	TheatreID   uuid.UUID
	Address     string
	TheatreName string
	TheatreArea string
}

func NewShowRepository(db *gorm.DB) *ShowRepository {
	return &ShowRepository{
		db: db,
	}
}

func (r *ShowRepository) Create(show *models.Show) (uuid.UUID, error) {
	err := r.db.Create(show).Error
	if err != nil {
		return uuid.Nil, err
	}
	return show.ID, err
}

func (r *ShowRepository) FindByMovie(movieid string) ([]models.Show, error) {
	var show []models.Show
	err := r.db.Where("movie_id = ?", movieid).Find(&show).Error
	return show, err
}

func (r *ShowRepository) FindById(id string) ([]models.Show, error) {
	var show []models.Show
	err := r.db.Where("id = ?", id).First(&show).Error
	return show, err
}

func (r *ShowRepository) Update(show *models.Show) error {
	return r.db.Save(show).Error
}

func (r *ShowRepository) FindByMovieAndCity(movieid string, city string, date string) ([]ShowTheatre, error) {
	var result []ShowTheatre

	err := r.db.
		Table("shows").
		Select(`
			shows.id as show_id,
			shows.start_time,
			shows.end_time,
			shows.regular_price as base_price,
			screens.screen_name as screen_name,
			shows.date,
			theatres.id as theatre_id,
			theatres.address,
			theatres.name as theatre_name,
			theatres.area as theatre_area
		`).
		Joins("JOIN screens ON screens.id = shows.screen_id").
		Joins("JOIN theatres ON theatres.id = screens.theatre_id").
		Where("shows.movie_id = ? AND theatres.city = ? AND shows.date = ? AND shows.start_time > ?", movieid, city, date, time.Now()).
		Scan(&result).Error

	return result, err
}

func (r *ShowRepository) CreateBookSeat(showID uuid.UUID) error {
	var show models.Show
	if err := r.db.First(&show, "id = ?", showID).Error; err != nil {
		return err
	}

	var seats []models.Seat
	if err := r.db.Where("screen_id = ?", show.ScreenID).Find(&seats).Error; err != nil {
		return err
	}

	var showSeats []models.ShowSeat

	for _, seat := range seats {
		price := 0.0

		switch seat.SeatType {
		case models.SeatRegular:
			price = show.RegularPrice
		case models.SeatPremium:
			price = show.PremiumPrice
		case models.SeatVIP:
			price = show.VIPPrice
		}

		showSeats = append(showSeats, models.ShowSeat{
			ID:         uuid.New(),
			ShowID:     show.ID,
			SeatID:     seat.ID,
			SeatNumber: seat.SeatNumber,
			SeatType:   string(seat.SeatType),
			SeatPrice:  price,
			Available:  true,
		})
	}

	return r.db.Create(&showSeats).Error
}

func (r *ShowRepository) CreateShowsFor15Days(show *models.Show) error {
	for i := 0; i < 15; i++ {
		newShow := *show
		newShow.ID = uuid.New()
		newShow.Date = show.Date.AddDate(0, 0, i)
		newShow.StartTime = show.StartTime.AddDate(0, 0, i)
		newShow.EndTime = show.EndTime.AddDate(0, 0, i)
		if err := r.db.Create(&newShow).Error; err != nil {
			return err
		}
		if err := r.CreateBookSeat(newShow.ID); err != nil {
			return err
		}
	}
	return nil
}

func (r *BookingRepository) GetShowStartTime(showID uuid.UUID) (time.Time, error) {
	var show models.Show
	err := r.DB.Select("start_time").First(&show, "id = ?", showID).Error
	if err != nil {
		return time.Time{}, err
	}
	return show.StartTime, nil
}
