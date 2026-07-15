package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/handlers"
	"github.com/poorna-9/goshow/internal/middleware"
)

func RegisterAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler, jwtSecret string) {
	auth := router.Group("/auth")
	auth.POST("/signup", authHandler.Signup)
	auth.POST("/login", authHandler.Login)
	auth.GET("/me", middleware.RequireAuth(jwtSecret), authHandler.Me)
}
