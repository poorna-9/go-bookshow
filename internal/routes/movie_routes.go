package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
	"github.com/poorna-9/goshow/internal/middleware"
)

func RegisterMovieRoutes(router *gin.RouterGroup, movieHandler *handlers.MovieHandler, jwtSecret string) {
	movies := router.Group("/movies")
	movies.POST("", middleware.RequireAuth(jwtSecret), middleware.RequireAdmin(), movieHandler.CreateMovie)
	movies.GET("/all", movieHandler.GetAllMovies)
	movies.GET("", movieHandler.GetMovieByCity)
	movies.GET("/:id", movieHandler.GetMovieById)
	movies.POST("/:id/poster", middleware.RequireAuth(jwtSecret), middleware.RequireAdmin(), movieHandler.UploadPoster)
}
