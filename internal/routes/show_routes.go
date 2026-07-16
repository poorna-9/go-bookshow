package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
	"github.com/poorna-9/goshow/internal/middleware"
)

func RegisterShowRoutes(router *gin.RouterGroup, showHandler *handlers.ShowHandler, jwtSecret string) {
	shows := router.Group("/shows")
	shows.POST("", middleware.RequireAuth(jwtSecret), middleware.RequireAdmin(), showHandler.CreateShow)
	shows.GET("/:id", showHandler.GetShowById)
	shows.PUT("/:id", middleware.RequireAuth(jwtSecret), middleware.RequireAdmin(), showHandler.UpdateShow)

	movies := router.Group("/movies")
	movies.GET("/:id/shows", showHandler.GetShowByMovieCity)
}
