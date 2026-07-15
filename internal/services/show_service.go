package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/repositories"
)

type ShowService struct {
	repo *repositories.ShowRepository
}

type ShowSlot struct {
	ShowID     uuid.UUID
	StartTime  time.Time
	EndTime    time.Time
	BasePrice  float64
	ScreenName string
}

type TheatreShows struct {
	TheatreID      uuid.UUID
	TheatreName    string
	TheatreArea    string
	TheatreAddress string
	Date           string
	Shows          []ShowSlot
}

func NewShowService(repo *repositories.ShowRepository) *ShowService {
	return &ShowService{repo: repo}
}

func (s *ShowService) CreateShow(show *models.Show) error {
	if show.MovieID == uuid.Nil {
		return errors.New("movie_id is required")
	}
	if show.ScreenID == uuid.Nil {
		return errors.New("screen_id is required")
	}
	if show.StartTime.IsZero() || show.EndTime.IsZero() || !show.EndTime.After(show.StartTime) {
		return errors.New("end time must be after start time")
	}
	if show.PremiumPrice <= 0 || show.RegularPrice <= 0 || show.VIPPrice <= 0 {
		return errors.New("prices should be greater than zero")
	}
	show.Date = show.StartTime.Truncate(24 * time.Hour)
	err := s.repo.CreateShowsFor15Days(show)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShowService) GetShowsByMovie(movieid string) ([]models.Show, error) {
	return s.repo.FindByMovie(movieid)
}

func (s *ShowService) GetShowById(id string) ([]models.Show, error) {
	return s.repo.FindById(id)
}

func (s *ShowService) UpdateShow(id string, updated *models.Show) error {
	shows, err := s.repo.FindById(id)
	if err != nil || len(shows) == 0 {
		return errors.New("show not found")
	}

	show := shows[0]
	show.MovieID = updated.MovieID
	show.ScreenID = updated.ScreenID
	show.StartTime = updated.StartTime
	show.EndTime = updated.EndTime
	show.RegularPrice = updated.RegularPrice
	show.PremiumPrice = updated.PremiumPrice
	show.VIPPrice = updated.PremiumPrice
	show.Date = updated.Date

	if show.StartTime.IsZero() || show.EndTime.IsZero() || !show.EndTime.After(show.StartTime) {
		return errors.New("timings are not valid")
	}
	if show.PremiumPrice <= 0 || show.RegularPrice <= 0 || show.VIPPrice <= 0 {
		return errors.New("prices should be greater than zero")
	}
	return s.repo.Update(&show)
}

func (s *ShowService) GetByMovieAndCity(movieid string, city string, date string) ([]TheatreShows, error) {
	flatrows, err := s.repo.FindByMovieAndCity(movieid, city, date)
	if err != nil {
		return nil, err
	}

	grouped := make(map[uuid.UUID]*TheatreShows)

	for _, row := range flatrows {
		theatre, exists := grouped[row.TheatreID]
		if !exists {
			theatre = &TheatreShows{
				TheatreID:      row.TheatreID,
				TheatreName:    row.TheatreName,
				TheatreArea:    row.TheatreArea,
				TheatreAddress: row.Address,
				Date:           date,
				Shows:          []ShowSlot{},
			}
			grouped[row.TheatreID] = theatre
		}
		theatre.Shows = append(theatre.Shows, ShowSlot{
			ShowID:     row.ShowID,
			StartTime:  row.StartTime,
			EndTime:    row.EndTime,
			BasePrice:  row.BasePrice,
			ScreenName: row.ScreenName,
		})
	}
	var result []TheatreShows
	for _, theatre := range grouped {
		result = append(result, *theatre)
	}

	return result, nil
}
