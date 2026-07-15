package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
)

func RegisterSeatRoutes(router *gin.RouterGroup, seatHandler *handlers.SeatHandler) {
	screens := router.Group("/screens")
	screens.GET("/:screen_id/seats", seatHandler.GetSeatsByScreen)
}
