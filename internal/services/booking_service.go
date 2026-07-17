package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/repositories"
)

const SessionTTL = 15 * time.Minute
const MaxSlotsPerSession = 5

type BookingService struct {
	Repo          *repositories.BookingRepository
	RazorpayKeyID string
}

func NewBookingService(repo *repositories.BookingRepository, razorpayKeyID string) *BookingService {
	return &BookingService{Repo: repo, RazorpayKeyID: razorpayKeyID}
}

type ReserveSlotResult struct {
	Action    string    `json:"action"`
	SessionID uuid.UUID `json:"session_id"`
}

type CheckoutResult struct {
	Amount          float64 `json:"amount"`
	RazorpayOrderID string  `json:"razorpay_order_id"`
	RazorpayKeyID   string  `json:"razorpay_key_id"`
	Expired         bool    `json:"expired"`
}

func (s *BookingService) ReserveSlot(userID, seatID, showID uuid.UUID) (*ReserveSlotResult, error) {

	activeShowStr, err := s.Repo.GetActivesession_in_redis(userID)
	startTime, err := s.Repo.GetShowStartTime(showID)
	if err != nil {
		return nil, errors.New("show not found")
	}
	if time.Now().After(startTime) {
		return nil, errors.New("this show has already started")
	}
	if activeShowStr != "" && activeShowStr != showID.String() {
		oldShowID, parseErr := uuid.Parse(activeShowStr)
		if parseErr == nil {
			if err := s.Repo.ReleaseUserSession(userID, oldShowID); err != nil {
				return nil, err
			}
		}
	}
	if err := s.Repo.SetActivityShow(userID, showID, SessionTTL); err != nil {
		return nil, err
	}
	session, err := s.Repo.GetActiveSession(userID, showID)
	if err != nil {
		session, err = s.Repo.CreateActiveSession(userID, showID)
		if err != nil {
			return nil, err
		}
	}

	alreadySelected, err := s.Repo.IsSeatInSession(session.ID, seatID)
	if err != nil {
		return nil, err
	}

	if alreadySelected {
		if err := s.Repo.RedisUnblock(session.ID, showID, seatID); err != nil {
			return nil, err
		}
		return &ReserveSlotResult{Action: "unselected", SessionID: session.ID}, nil
	}

	count, err := s.Repo.SessionSeatCount(session.ID)
	if err != nil {
		return nil, err
	}
	if count >= MaxSlotsPerSession {
		return nil, errors.New("max seats per session reached")
	}

	if err := s.Repo.RedisBlock(session.ID, showID, seatID, SessionTTL); err != nil {
		return nil, err
	}

	return &ReserveSlotResult{Action: "selected", SessionID: session.ID}, nil
}

func (s *BookingService) GetSeatsForShow(showID uuid.UUID) ([]models.ShowSeat, error) {
	return s.Repo.GetSeatsOfShow(showID)
}

type SeatStatus struct {
	SeatID     uuid.UUID `json:"seat_id"`
	SeatNumber string    `json:"seat_number"`
	SeatType   string    `json:"seat_type"`
	Price      float64   `json:"price"`
	Status     string    `json:"status"` // "available" | "yours" | "taken" | "booked"
}

type ReservedSlotsResult struct {
	Seats     []SeatStatus `json:"seats"`
	SessionID *uuid.UUID   `json:"session_id"`
}

func (s *BookingService) GetReservedSlots(userID *uuid.UUID, showID uuid.UUID) (*ReservedSlotsResult, error) {
	showSeats, err := s.Repo.GetSeatsOfShow(showID)
	if err != nil {
		return nil, err
	}

	result := &ReservedSlotsResult{Seats: []SeatStatus{}}

	var userSessionID *uuid.UUID
	if userID != nil {
		session, err := s.Repo.GetActiveSession(*userID, showID)
		if err == nil {
			userSessionID = &session.ID
			result.SessionID = &session.ID
		}
	}

	for _, seat := range showSeats {
		status := "available"

		if !seat.Available {
			status = "booked"
		} else {
			lockedBySession, err := s.Repo.GetLockSession(showID, seat.SeatID)
			if err != nil {
				return nil, err
			}
			if lockedBySession != "" {
				if userSessionID != nil && lockedBySession == userSessionID.String() {
					status = "yours"
				} else {
					status = "taken"
				}
			}
		}

		result.Seats = append(result.Seats, SeatStatus{
			SeatID:     seat.SeatID,
			SeatNumber: seat.SeatNumber,
			SeatType:   seat.SeatType,
			Price:      seat.SeatPrice,
			Status:     status,
		})
	}

	return result, nil
}

type CheckoutSummary struct {
	SessionID   uuid.UUID  `json:"session_id"`
	ShowID      uuid.UUID  `json:"show_id"`
	Seats       []SeatItem `json:"seats"`
	TotalAmount float64    `json:"total_amount"`
}

type SeatItem struct {
	SeatID     uuid.UUID `json:"seat_id"`
	SeatNumber string    `json:"seat_number"`
	SeatType   string    `json:"seat_type"`
	Price      float64   `json:"price"`
}

func (s *BookingService) GetCheckoutSummary(userid uuid.UUID, showid uuid.UUID) (*CheckoutSummary, error) {
	session, err := s.Repo.GetActiveSession(userid, showid)
	if err != nil {
		return nil, errors.New("no active session found")
	}
	seatIDs, err := s.Repo.GetSessionSeatIDs(session.ID)
	if err != nil {
		return nil, err
	}
	if len(seatIDs) == 0 {
		return nil, errors.New("no seats selected")
	}

	showSeats, err := s.Repo.GetSeatsByIds(showid, seatIDs)
	if err != nil {
		return nil, err
	}
	var items []SeatItem
	var total float64
	for _, seat := range showSeats {
		items = append(items, SeatItem{
			SeatID:     seat.SeatID,
			SeatNumber: seat.SeatNumber,
			SeatType:   seat.SeatType,
			Price:      seat.SeatPrice,
		})
		total += seat.SeatPrice
	}
	return &CheckoutSummary{
		SessionID:   session.ID,
		ShowID:      showid,
		Seats:       items,
		TotalAmount: total,
	}, nil
}

func (s *BookingService) InitiateCheckOut(userID, showID uuid.UUID) (*CheckoutResult, error) {
	session, err := s.Repo.GetActiveSession(userID, showID)
	if err != nil {
		return nil, errors.New("no active session found")
	}

	existingPayment, err := s.Repo.GetPendingPaymentBySession(session.ID)
	if err == nil {
		remaining := time.Until(existingPayment.ExpiresAt)
		if remaining > 0 {
			return &CheckoutResult{
				Amount:          existingPayment.Amount,
				RazorpayOrderID: existingPayment.RazorpayOrderID,
				RazorpayKeyID:   s.RazorpayKeyID,
			}, nil
		}

		seatIDs, _ := s.Repo.GetSessionSeatIDs(session.ID)
		if err := s.Repo.ExpireSessionTx(session.ID, showID, seatIDs); err != nil {
			return nil, err
		}
		return &CheckoutResult{Expired: true}, nil
	}

	seatIDs, err := s.Repo.GetSessionSeatIDs(session.ID)
	if err != nil {
		return nil, err
	}
	if len(seatIDs) == 0 {
		return nil, errors.New("no seats selected")
	}

	showSeats, err := s.Repo.GetSeatsByIds(showID, seatIDs)
	if err != nil {
		return nil, err
	}

	var total float64
	for _, seat := range showSeats {
		total += seat.SeatPrice
	}

	if err := s.Repo.BlockSeatsAndSnapshotTx(session.ID, showID, seatIDs); err != nil {
		return nil, err
	}

	orderID, err := s.Repo.CreateRazorpayOrder(total, session.ID.String())
	if err != nil {
		return nil, err
	}

	payment := &models.Payment{
		ID:              uuid.New(),
		UserID:          userID,
		SessionID:       session.ID,
		ShowID:          showID,
		Amount:          total,
		Status:          models.PaymentPending,
		Method:          "razorpay",
		RazorpayOrderID: orderID,
		ExpiresAt:       time.Now().Add(10 * time.Minute),
	}
	if err := s.Repo.CreatePayment(payment); err != nil {
		return nil, err
	}

	return &CheckoutResult{
		Amount:          total,
		RazorpayOrderID: orderID,
		RazorpayKeyID:   s.RazorpayKeyID,
	}, nil
}

func (s *BookingService) FinalizePayment(orderID, razorpayPaymentID string) (*models.Booking, error) {
	payment, err := s.Repo.GetPaymentByOrderID(orderID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	if payment.Status == models.PaymentSuccess {
		return s.Repo.GetBookingBySessionID(payment.SessionID)
	}

	return s.Repo.Finalizetx(payment, razorpayPaymentID)
}

func (s *BookingService) VerifySignature(orderID, paymentID, signature string) bool {
	return s.Repo.VerifyPaymentSignature(orderID, paymentID, signature)
}

func (s *BookingService) VerifyWebhookSignature(payload []byte, signature string) bool {
	return s.Repo.VerifyWebhookSignature(payload, signature)
}

type CancelResult struct {
	CanRetry bool `json:"can_retry"`
}

func (s *BookingService) HandlePaymentCancel(orderID string) (*CancelResult, error) {
	payment, err := s.Repo.GetPaymentByOrderID(orderID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	if payment.Status != models.PaymentPending {
		return &CancelResult{CanRetry: false}, nil
	}

	remaining := time.Until(payment.ExpiresAt)
	if remaining > 0 {
		return &CancelResult{CanRetry: true}, nil
	}

	if err := s.Repo.CancelPaymentTx(payment); err != nil {
		return nil, err
	}
	return &CancelResult{CanRetry: false}, nil
}

type PaymentStatusResult struct {
	Status   string `json:"status"`
	Redirect string `json:"redirect,omitempty"`
}

func (s *BookingService) GetPaymentStatus(orderID string) (*PaymentStatusResult, error) {
	payment, err := s.Repo.GetPaymentByOrderID(orderID)
	if err != nil {
		return &PaymentStatusResult{Status: "failed", Redirect: "/movies"}, nil
	}

	switch payment.Status {
	case models.PaymentSuccess:
		booking, err := s.Repo.GetBookingBySessionID(payment.SessionID)
		if err != nil {
			return &PaymentStatusResult{Status: "pending"}, nil
		}
		return &PaymentStatusResult{
			Status:   "success",
			Redirect: fmt.Sprintf("booking.html?booking_id=%s", booking.ID),
		}, nil

	case models.PaymentFailed:
		return &PaymentStatusResult{Status: "failed", Redirect: "/movies"}, nil

	default: // pending
		remaining := time.Until(payment.ExpiresAt)
		if remaining <= 0 {
			_, _ = s.HandlePaymentCancel(orderID)
			return &PaymentStatusResult{Status: "expired", Redirect: "/movies"}, nil
		}
		return &PaymentStatusResult{Status: "pending"}, nil
	}
}

type BookingDetail struct {
	BookingID   uuid.UUID  `json:"booking_id"`
	Status      string     `json:"status"`
	TotalAmount float64    `json:"total_amount"`
	ShowID      uuid.UUID  `json:"show_id"`
	Seats       []SeatItem `json:"seats"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (s *BookingService) GetBookingDetail(bookingID, userID uuid.UUID) (*BookingDetail, error) {
	booking, err := s.Repo.GetBookingByID(bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}
	if booking.UserID != userID {
		return nil, errors.New("not authorized to view this booking")
	}

	showSeats, err := s.Repo.GetSeatsByIds(booking.ShowID, booking.SeatIDs)
	if err != nil {
		return nil, err
	}

	var items []SeatItem
	for _, seat := range showSeats {
		items = append(items, SeatItem{
			SeatID:     seat.SeatID,
			SeatNumber: seat.SeatNumber,
			SeatType:   seat.SeatType,
			Price:      seat.SeatPrice,
		})
	}

	return &BookingDetail{
		BookingID:   booking.ID,
		Status:      string(booking.Status),
		TotalAmount: booking.TotalAmount,
		ShowID:      booking.ShowID,
		Seats:       items,
		CreatedAt:   booking.CreatedAt,
	}, nil
}

func (s *BookingService) SweepStaleBookings() {
	cutoff := time.Now().Add(-SessionTTL)
	staleSessions, err := s.Repo.FindStaleSessions(cutoff)
	if err == nil {
		for _, session := range staleSessions {
			seatIDs, _ := s.Repo.GetSessionSeatIDs(session.ID)
			for _, seatID := range seatIDs {
				s.Repo.RedisUnblock(session.ID, session.ShowID, seatID)
			}
			s.Repo.ExpireSessionTx(session.ID, session.ShowID, seatIDs)
		}
	}

	stalePayments, err := s.Repo.FindStalePendingPayments()
	if err == nil {
		for _, payment := range stalePayments {
			s.HandlePaymentCancel(payment.RazorpayOrderID)
		}
	}
}
