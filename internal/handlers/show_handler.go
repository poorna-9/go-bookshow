package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/services"
)

type ShowHandler struct {
	showservice *services.ShowService
}

func NewShowHandler(showservice *services.ShowService) *ShowHandler {
	return &ShowHandler{showservice: showservice}
}

func (h *ShowHandler) CreateShow(c *gin.Context) {
	var show models.Show
	err := c.ShouldBindJSON(&show)
	if err != nil {
		c.JSON(500, gin.H{"context": err.Error()})
		return
	}
	newerr := h.showservice.CreateShow(&show)
	if newerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": newerr.Error()})
		return
	}
	c.JSON(200, show)
}

func (h *ShowHandler) UpdateShow(c *gin.Context) {
	id := c.Param("id")
	var updated models.Show
	err := c.ShouldBindBodyWithJSON(&updated)
	if err != nil {
		c.JSON(500, gin.H{"context": err.Error()})
		return
	}
	newerr := h.showservice.UpdateShow(id, &updated)
	if newerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": newerr.Error()})
		return
	}
	c.JSON(200, updated)
}

func (h *ShowHandler) GetShowById(c *gin.Context) {
	id := c.Param("id")
	show, err := h.showservice.GetShowById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, show)
}

func (h *ShowHandler) GetShowByMovie(c *gin.Context) {
	movieid := c.Param("movie_id")
	show, err := h.showservice.GetShowsByMovie(movieid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, show)
}

func (h *ShowHandler) GetShowByMovieCity(c *gin.Context) {
	movieid := c.Param("movie_id")
	city := c.Query("city")
	date := c.Query("date")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city query parameter is required"})
		return
	}
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date query parameter is required"})
		return
	}

	result, err := h.showservice.GetByMovieAndCity(movieid, city, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
