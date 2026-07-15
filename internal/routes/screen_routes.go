package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
)

func RegisterScreenRoutes(router *gin.RouterGroup, screenHandler *handlers.ScreenHandler) {
	screens := router.Group("/screens")
	screens.POST("", screenHandler.CreateScreen)
	screens.GET("/:id", screenHandler.GetScreenInfo)
	screens.DELETE("/:id", screenHandler.DeactivateScreen)

	theatres := router.Group("/theatres")
	theatres.GET("/:theatre_id/screens", screenHandler.GetScreensByTheatre)
}