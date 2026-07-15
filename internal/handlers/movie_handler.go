package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/services"
)

type MovieHandler struct {
	Service *services.MovieService
}

func NewMovieHandler(service *services.MovieService) *MovieHandler {
	return &MovieHandler{Service: service}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var movie models.Movie
	err := c.ShouldBindBodyWithJSON(&movie)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newerr := h.Service.CreateMovie(&movie)
	if newerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": newerr.Error()})
		return
	}
	c.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) GetMovieById(c *gin.Context) {
	id := c.Param("id")
	movie, err := h.Service.GetMovieById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movie)

}

func (h *MovieHandler) GetMovieByCity(c *gin.Context) {
	city := c.Param("city")
	movies, err := h.Service.GetMovieByCity(city)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(200, movies)
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.Service.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}
