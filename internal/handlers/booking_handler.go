package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/poorna-9/goshow/internal/services"
)

type BookingHandler struct {
	Service *services.BookingService
}

func NewBookingHandler(service *services.BookingService) *BookingHandler {
	return &BookingHandler{Service: service}
}

type ReserveSlotInput struct {
	UserID uuid.UUID `json:"user_id"`
	ShowID uuid.UUID `json:"show_id"`
	SeatID uuid.UUID `json:"seat_id"`
}

func (h *BookingHandler) ReserveSlot(c *gin.Context) {
	var input ReserveSlotInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.Service.ReserveSlot(input.UserID, input.SeatID, input.ShowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *BookingHandler) GetReservedSlots(c *gin.Context) {
	showID, err := uuid.Parse(c.Param("show_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid show_id"})
		return
	}

	var userID *uuid.UUID
	if uidStr := c.Query("user_id"); uidStr != "" {
		parsed, err := uuid.Parse(uidStr)
		if err == nil {
			userID = &parsed
		}
	}

	result, err := h.Service.GetReservedSlots(userID, showID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

type CheckoutInput struct {
	UserID uuid.UUID `json:"user_id"`
	ShowID uuid.UUID `json:"show_id"`
}

func (h *BookingHandler) GetCheckout(c *gin.Context) {
	showID, err := uuid.Parse(c.Query("show_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid show_id"})
		return
	}
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	summary, err := h.Service.GetCheckoutSummary(userID, showID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *BookingHandler) InitiateCheckout(c *gin.Context) {
	var input CheckoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.Service.InitiateCheckOut(input.UserID, input.ShowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

type PaymentCallbackInput struct {
	RazorpayOrderID   string `json:"razorpay_order_id"`
	RazorpayPaymentID string `json:"razorpay_payment_id"`
	RazorpaySignature string `json:"razorpay_signature"`
	Error             bool   `json:"error"`
}

func (h *BookingHandler) PaymentCallback(c *gin.Context) {
	var input PaymentCallbackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Error {
		result, err := h.Service.HandlePaymentCancel(input.RazorpayOrderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": false, "can_retry": result.CanRetry})
		return
	}

	valid := h.Service.VerifySignature(input.RazorpayOrderID, input.RazorpayPaymentID, input.RazorpaySignature)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid signature"})
		return
	}

	booking, err := h.Service.FinalizePayment(input.RazorpayOrderID, input.RazorpayPaymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "booking_id": booking.ID})
}

func (h *BookingHandler) RazorpayWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	signature := c.GetHeader("X-Razorpay-Signature")
	if !h.Service.VerifyWebhookSignature(body, signature) {
		c.Status(http.StatusBadRequest)
		return
	}

	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	eventType, _ := event["event"].(string)
	payload, _ := event["payload"].(map[string]interface{})
	paymentWrapper, _ := payload["payment"].(map[string]interface{})
	paymentEntity, _ := paymentWrapper["entity"].(map[string]interface{})

	orderID, _ := paymentEntity["order_id"].(string)
	paymentID, _ := paymentEntity["id"].(string)

	switch eventType {
	case "payment.captured":
		h.Service.FinalizePayment(orderID, paymentID)
	case "payment.failed":
		h.Service.HandlePaymentCancel(orderID)
	}

	c.Status(http.StatusOK)
}
