package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
)

func RegisterMovieRoutes(router *gin.RouterGroup, movieHandler *handlers.MovieHandler) {
	movies := router.Group("/movies")
	movies.POST("", movieHandler.CreateMovie)
	movies.GET("", movieHandler.GetAllMovies)
	movies.GET("/:id", movieHandler.GetMovieById)
	movies.GET("", movieHandler.GetMovieByCity)
}
