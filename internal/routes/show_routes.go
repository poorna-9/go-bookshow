package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
)

func RegisterShowRoutes(router *gin.RouterGroup, showHandler *handlers.ShowHandler) {
	shows := router.Group("/shows")
	shows.POST("", showHandler.CreateShow)
	shows.GET("/:id", showHandler.GetShowById)
	shows.PUT("/:id", showHandler.UpdateShow)

	movies := router.Group("/movies")
	movies.GET("/:movie_id/shows", showHandler.GetShowByMovieCity)
}
