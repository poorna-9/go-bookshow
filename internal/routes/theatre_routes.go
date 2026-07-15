package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
)

func RegisterTheatreRoutes(router *gin.RouterGroup, theatreHandler *handlers.TheatreHandler) {
	theatres := router.Group("/theatres")
	theatres.POST("", theatreHandler.CreateTheatre)
	theatres.GET("", theatreHandler.GetTheatresByCity)
	theatres.GET("/:id", theatreHandler.GetTheatreInfo)
	theatres.GET("/all", theatreHandler.GetAllTheatres)
}
