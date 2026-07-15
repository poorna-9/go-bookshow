package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/services"
)

type TheatreHandler struct {
	service *services.TheatreService
}

func NewTheatreHandler(service *services.TheatreService) *TheatreHandler {
	return &TheatreHandler{
		service: service,
	}
}

func (h *TheatreHandler) CreateTheatre(c *gin.Context) {
	var theatre models.Theatre
	c.ShouldBindJSON(&theatre)
	createdTheatre := h.service.CreateTheatre(&theatre)
	c.JSON(201, gin.H{
		"message": "theatre created successfully",
		"data":    createdTheatre,
	})
}

func (h *TheatreHandler) GetTheatresByCity(c *gin.Context) {
	city := c.Query("city")
	theatres := h.service.GetTheatresByCity(city)
	c.JSON(200, gin.H{
		"data": theatres,
	})
}

func (h *TheatreHandler) GetTheatreInfo(c *gin.Context) {
	id := c.Param("id")
	theatre := h.service.Gettheatreinfo(id)
	c.JSON(200, gin.H{
		"data": theatre,
	})
}

func (h *TheatreHandler) GetAllTheatres(c *gin.Context) {
	theatres := h.service.GetAllTheatres()
	c.JSON(http.StatusOK, theatres)
}
