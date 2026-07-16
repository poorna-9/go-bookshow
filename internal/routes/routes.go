package routes

import (
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

func RegisterRoutes(router *gin.Engine, h *Handlers, jwtSecret string) {
	v1 := router.Group("/api/v1")

	RegisterTheatreRoutes(v1, h.Theatre, jwtSecret)
	RegisterScreenRoutes(v1, h.Screen, jwtSecret)
	RegisterSeatRoutes(v1, h.Seat)
	RegisterMovieRoutes(v1, h.Movie, jwtSecret)
	RegisterShowRoutes(v1, h.Show, jwtSecret)
	RegisterBookingRoutes(v1, h.Booking, jwtSecret)
	RegisterAuthRoutes(v1, h.Auth, jwtSecret)
}
