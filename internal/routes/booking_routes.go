package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
)

func RegisterBookingRoutes(router *gin.RouterGroup, bookingHandler *handlers.BookingHandler) {
	bookings := router.Group("/bookings")
	bookings.POST("/reserve", bookingHandler.ReserveSlot)
	bookings.GET("/checkout", bookingHandler.GetCheckout)
	bookings.POST("/checkout", bookingHandler.InitiateCheckout)
	bookings.POST("/payment-callback", bookingHandler.PaymentCallback)
	bookings.POST("/webhook/razorpay", bookingHandler.RazorpayWebhook)
	bookings.GET("/payment-status/:order_id", bookingHandler.GetPaymentStatus)

	shows := router.Group("/shows")
	shows.GET("/:show_id/reserved-slots", bookingHandler.GetReservedSlots)
}
