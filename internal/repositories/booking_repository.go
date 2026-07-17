package repositories

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/razorpay/razorpay-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/poorna-9/goshow/internal/models"
)

type BookingRepository struct {
	DB                    *gorm.DB
	RDB                   *redis.Client
	RZP                   *razorpay.Client
	RazorpaySecret        string
	RazorpayWebhookSecret string
}

func NewBookingRepository(db *gorm.DB, rdb *redis.Client, rzp *razorpay.Client, razorpaySecret, razorpayWebhookSecret string) *BookingRepository {
	return &BookingRepository{
		DB:                    db,
		RDB:                   rdb,
		RZP:                   rzp,
		RazorpaySecret:        razorpaySecret,
		RazorpayWebhookSecret: razorpayWebhookSecret,
	}
}

func (r *BookingRepository) RedisBlock(bookingSession, showID, seatID uuid.UUID, ttl time.Duration) error {
	ctx := context.Background()
	sessionKey := fmt.Sprintf("session_seats:%s", bookingSession.String())
	lockKey := fmt.Sprintf("seat_lock:%s:%s", showID.String(), seatID.String())

	locked, err := r.RDB.SetNX(ctx, lockKey, bookingSession.String(), ttl).Result()
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("seat already selected")
	}

	if err := r.RDB.SAdd(ctx, sessionKey, seatID.String()).Err(); err != nil {
		r.RDB.Del(ctx, lockKey)
		return err
	}
	return nil
}

func (r *BookingRepository) RedisUnblock(bookingSession, showID, seatID uuid.UUID) error {
	ctx := context.Background()
	sessionKey := fmt.Sprintf("session_seats:%s", bookingSession.String())
	lockKey := fmt.Sprintf("seat_lock:%s:%s", showID.String(), seatID.String())

	if err := r.RDB.SRem(ctx, sessionKey, seatID.String()).Err(); err != nil {
		return err
	}
	return r.RDB.Del(ctx, lockKey).Err()
}

func (r *BookingRepository) GetLockSession(showID, seatID uuid.UUID) (string, error) {
	ctx := context.Background()
	lockKey := fmt.Sprintf("seat_lock:%s:%s", showID.String(), seatID.String())

	val, err := r.RDB.Get(ctx, lockKey).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *BookingRepository) IsSeatInSession(bookingSession, seatID uuid.UUID) (bool, error) {
	ctx := context.Background()
	sessionKey := fmt.Sprintf("session_seats:%s", bookingSession.String())
	return r.RDB.SIsMember(ctx, sessionKey, seatID.String()).Result()
}

func (r *BookingRepository) SessionSeatCount(bookingSession uuid.UUID) (int64, error) {
	ctx := context.Background()
	sessionKey := fmt.Sprintf("session_seats:%s", bookingSession.String())
	return r.RDB.SCard(ctx, sessionKey).Result()
}

func (r *BookingRepository) GetActiveSession(userID, showID uuid.UUID) (*models.BookingSession, error) {
	var session models.BookingSession
	err := r.DB.Where("user_id = ? AND show_id = ? AND expired = ? AND success = ?", userID, showID, false, false).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *BookingRepository) CreateActiveSession(userID, showID uuid.UUID) (*models.BookingSession, error) {
	session := models.BookingSession{
		ID:      uuid.New(),
		UserID:  userID,
		ShowID:  showID,
		Expired: false,
		Success: false,
	}
	err := r.DB.Create(&session).Error
	return &session, err
}

func (r *BookingRepository) GetSeatsOfShow(showID uuid.UUID) ([]models.ShowSeat, error) {
	var seats []models.ShowSeat
	err := r.DB.Where("show_id = ?", showID).Find(&seats).Error
	return seats, err
}

func (r *BookingRepository) GetSessionSeatIDs(bookingSession uuid.UUID) ([]uuid.UUID, error) {
	ctx := context.Background()
	sessionKey := fmt.Sprintf("session_seats:%s", bookingSession.String())

	rawIDs, err := r.RDB.SMembers(ctx, sessionKey).Result()
	if err != nil {
		return nil, err
	}

	var seatIDs []uuid.UUID
	for _, raw := range rawIDs {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			continue
		}
		seatIDs = append(seatIDs, parsed)
	}

	return seatIDs, nil
}

func (r *BookingRepository) GetSeatsByIds(showid uuid.UUID, seatids []uuid.UUID) ([]models.ShowSeat, error) {
	var seats []models.ShowSeat
	err := r.DB.Where("show_id = ? AND seat_id IN ?", showid, seatids).Find(&seats).Error
	if err != nil {
		return nil, err
	}
	return seats, err
}

func (r *BookingRepository) CreateRazorpayOrder(amount float64, receipt string) (string, error) {
	data := map[string]interface{}{
		"amount":   int(amount * 100),
		"currency": "INR",
		"receipt":  receipt,
	}
	order, err := r.RZP.Order.Create(data, nil)
	if err != nil {
		return "", err
	}
	orderID, ok := order["id"].(string)
	if !ok {
		return "", errors.New("unexpected response from razorpay")
	}
	return orderID, nil
}

func (r *BookingRepository) GetPendingPaymentBySession(sessionid uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.DB.
		Where("session_id = ? AND status = ?", sessionid, models.PaymentPending).
		Order("created_at DESC").
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *BookingRepository) BlockSeatsAndSnapshotTx(session_id uuid.UUID, show_id uuid.UUID, seat_ids []uuid.UUID) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.ShowSeat{}).Where("show_id = ? AND seat_id IN ?", show_id, seat_ids).Update("is_blocked", true).Error
		if err != nil {
			return err
		}
		newerr := tx.Model(&models.BookingSession{}).Where("id = ?", session_id).Update("session_seats", seat_ids).Error
		if newerr != nil {
			return newerr
		}
		return nil
	})
}

func (r *BookingRepository) CreatePayment(payment *models.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *BookingRepository) ExpireSessionTx(session_id uuid.UUID, show_id uuid.UUID, seats []uuid.UUID) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.BookingSession{}).Where("id = ?", session_id).Update("expired", true).Error; err != nil {
			return err
		}
		if len(seats) > 0 {
			if err := tx.Model(&models.ShowSeat{}).
				Where("show_id = ? AND seat_id IN ?", show_id, seats).
				Update("is_blocked", false).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *BookingRepository) GetPaymentByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.DB.Where("razorpay_order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *BookingRepository) Finalizetx(payment *models.Payment, razorpayPaymentID string) (*models.Booking, error) {
	var booking models.Booking

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var session models.BookingSession
		newerr := r.DB.Where("id = ?", payment.SessionID).First(&session).Error
		if newerr != nil {
			return newerr
		}
		booking = models.Booking{
			ID:          uuid.New(),
			UserID:      payment.UserID,
			ShowID:      payment.ShowID,
			SessionID:   session.ID,
			Status:      models.BookingConfirmed,
			TotalAmount: payment.Amount,
			SeatIDs:     session.SessionSeats,
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.ShowSeat{}).
			Where("show_id = ? AND seat_id IN ?", payment.ShowID, session.SessionSeats).
			Updates(map[string]any{"available": false, "is_blocked": false}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Payment{}).
			Where("id = ?", payment.ID).
			Updates(map[string]any{"status": models.PaymentSuccess, "transaction_ref": razorpayPaymentID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.BookingSession{}).
			Where("id = ?", session.ID).
			Update("success", true).Error; err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	return &booking, err
}

func (r *BookingRepository) GetBookingBySessionID(sessionID uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.DB.Where("session_id = ?", sessionID).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) VerifyPaymentSignature(orderID, paymentID, signature string) bool {
	data := orderID + "|" + paymentID
	h := hmac.New(sha256.New, []byte(r.RazorpaySecret))
	h.Write([]byte(data))
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func (r *BookingRepository) CancelPaymentTx(payment *models.Payment) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Payment{}).
			Where("id = ?", payment.ID).
			Update("status", models.PaymentFailed).Error; err != nil {
			return err
		}

		var session models.BookingSession
		if err := tx.First(&session, "id = ?", payment.SessionID).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.ShowSeat{}).
			Where("show_id = ? AND seat_id IN ?", payment.ShowID, session.SessionSeats).
			Update("is_blocked", false).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.BookingSession{}).
			Where("id = ?", session.ID).
			Update("expired", true).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *BookingRepository) VerifyWebhookSignature(payload []byte, signature string) bool {
	h := hmac.New(sha256.New, []byte(r.RazorpayWebhookSecret))
	h.Write(payload)
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func (r *BookingRepository) GetPaymentStatus(orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.DB.Where("razorpay_order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *BookingRepository) SetActivityShow(user_id, show_id uuid.UUID, ttl time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("active_show:%s", user_id.String())
	return r.RDB.Set(ctx, key, show_id.String(), ttl).Err()
}
func (r *BookingRepository) GetActivesession_in_redis(user_id uuid.UUID) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("active_show:%s", user_id.String())
	val, err := r.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil

}

func (r *BookingRepository) ReleaseUserSession(user_id, show_id uuid.UUID) error {
	session, err := r.GetActiveSession(user_id, show_id)
	if err != nil {
		return nil
	}
	ctx := context.Background()
	seatIDs, err := r.GetSessionSeatIDs(session.ID)
	if err != nil {
		return err
	}

	for _, seatID := range seatIDs {
		lockKey := fmt.Sprintf("seat_lock:%s:%s", show_id.String(), seatID.String())
		r.RDB.Del(ctx, lockKey)
	}
	sessionKey := fmt.Sprintf("session_seats:%s", session.ID.String())
	r.RDB.Del(ctx, sessionKey)

	return r.DB.Model(&models.BookingSession{}).
		Where("id = ?", session.ID).
		Update("expired", true).Error

}

func (r *BookingRepository) GetBookingByID(id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.DB.First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) FindStaleSessions(older_than time.Time) ([]models.BookingSession, error) {
	var session []models.BookingSession
	err := r.DB.Where("expired = ? AND sucess = ? AND created_at < ?", false, false, older_than).Find(&session).Error
	return session, err
}

func (r *BookingRepository) FindStalePendingPayments() ([]models.Payment, error) {
	var payments []models.Payment
	err := r.DB.Where("status = ? AND expires_at < ?", models.PaymentPending, time.Now()).Find(&payments).Error
	return payments, err
}

func (r *BookingRepository) AreSeatsAvailable(showID uuid.UUID, seatIDs []uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.ShowSeat{}).
		Where("show_id = ? AND seat_id IN ? AND available = ?", showID, seatIDs, false).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *BookingRepository) GetSessionByID(sessionID uuid.UUID) (*models.BookingSession, error) {
	var session models.BookingSession
	err := r.DB.First(&session, "id = ?", sessionID).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *BookingRepository) MarkPaymentRefundRequired(paymentID uuid.UUID) error {
	return r.DB.Model(&models.Payment{}).
		Where("id = ?", paymentID).
		Update("status", models.PaymentRefundRequired).Error
}
