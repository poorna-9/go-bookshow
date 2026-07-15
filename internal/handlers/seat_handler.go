package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/poorna-9/goshow/internal/services"
)

type SeatHandler struct {
	Service *services.SeatService
}

func NewSeatHandler(service *services.SeatService) *SeatHandler {
	return &SeatHandler{Service: service}
}

func (h *SeatHandler) GetSeatsByScreen(c *gin.Context) {
	screenID := c.Param("screen_id")
	seats, err := h.Service.GetSeatsByScreen(screenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seats)
}
