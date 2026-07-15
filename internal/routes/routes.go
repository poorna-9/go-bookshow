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
}

func RegisterRoutes(router *gin.Engine, h *Handlers) {
	v1 := router.Group("/api/v1")

	RegisterTheatreRoutes(v1, h.Theatre)
	RegisterScreenRoutes(v1, h.Screen)
	RegisterSeatRoutes(v1, h.Seat)
	RegisterMovieRoutes(v1, h.Movie)
	RegisterShowRoutes(v1, h.Show)
	RegisterBookingRoutes(v1, h.Booking)
}
