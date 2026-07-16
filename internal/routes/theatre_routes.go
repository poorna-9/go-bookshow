package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
	"github.com/poorna-9/goshow/internal/middleware"
)

func RegisterTheatreRoutes(router *gin.RouterGroup, theatreHandler *handlers.TheatreHandler, jwtSecret string) {
	theatres := router.Group("/theatres")
	theatres.POST("", middleware.RequireAuth(jwtSecret), middleware.RequireAdmin(), theatreHandler.CreateTheatre)
	theatres.GET("/all", theatreHandler.GetAllTheatres)
	theatres.GET("", theatreHandler.GetTheatresByCity)
	theatres.GET("/:id", theatreHandler.GetTheatreInfo)
}
