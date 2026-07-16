package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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

func (h *MovieHandler) UploadPoster(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	file, err := c.FormFile("poster")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "poster file is required"})
		return
	}

	filename := id.String() + filepath.Ext(file.Filename)
	savePath := filepath.Join("web", "images", "posters", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save poster"})
		return
	}

	url := "/images/posters/" + filename
	if err := h.Service.UpdatePosterURL(id, url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"poster_url": url})
}
