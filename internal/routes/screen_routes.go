package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
	"github.com/poorna-9/goshow/internal/middleware"
)

func RegisterScreenRoutes(router *gin.RouterGroup, screenHandler *handlers.ScreenHandler, jwtSecret string) {
	screens := router.Group("/screens")
	screens.POST("", middleware.RequireAuth(jwtSecret), middleware.RequireAdmin(), screenHandler.CreateScreen)
	screens.GET("/:id", screenHandler.GetScreenInfo)
	screens.DELETE("/:id", middleware.RequireAuth(jwtSecret), middleware.RequireAdmin(), screenHandler.DeactivateScreen)

	theatres := router.Group("/theatres")
	theatres.GET("/:theatre_id/screens", screenHandler.GetScreensByTheatre)
}
