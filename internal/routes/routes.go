package routes

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/poorna-9/goshow/internal/handlers"
)

type Handlers struct {
	Theatre *handlers.TheatreHandler
	Screen  *handlers.ScreenHandler
	Seat    *handlers.SeatHandler
	Movie   *handlers.MovieHandler
	Show    *handlers.ShowHandler
	Booking *handlers.BookingHandler
	Auth    *handlers.AuthHandler
}

func RegisterRoutes(router *gin.Engine, h *Handlers) {
	v1 := router.Group("/api/v1")

	RegisterTheatreRoutes(v1, h.Theatre)
	RegisterScreenRoutes(v1, h.Screen)
	RegisterSeatRoutes(v1, h.Seat)
	RegisterMovieRoutes(v1, h.Movie)
	RegisterShowRoutes(v1, h.Show)
	RegisterBookingRoutes(v1, h.Booking)
	// jwtSecret is read from the JWT_SECRET env var with a default fallback
	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" {
		jwt = "secret"
	}
	RegisterAuthRoutes(v1, h.Auth, jwt)
}
