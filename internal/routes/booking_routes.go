package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
	"github.com/poorna-9/goshow/internal/middleware"
)

func RegisterBookingRoutes(router *gin.RouterGroup, bookingHandler *handlers.BookingHandler, jwtSecret string) {
	bookings := router.Group("/bookings")
	bookings.Use(middleware.RequireAuth(jwtSecret))
	bookings.POST("/reserve", bookingHandler.ReserveSlot)
	bookings.GET("/checkout", bookingHandler.GetCheckout)
	bookings.POST("/checkout", bookingHandler.InitiateCheckout)
	bookings.POST("/payment-callback", bookingHandler.PaymentCallback)
	bookings.GET("/payment-status/:order_id", bookingHandler.GetPaymentStatus)
	bookings.GET("/:id", bookingHandler.GetBookingDetail)

	// webhook stays PUBLIC — Razorpay calls this directly, no user token exists
	router.POST("/bookings/webhook/razorpay", bookingHandler.RazorpayWebhook)

	shows := router.Group("/shows")
	shows.Use(middleware.RequireAuth(jwtSecret))
	shows.GET("/:show_id/reserved-slots", bookingHandler.GetReservedSlots)
}
