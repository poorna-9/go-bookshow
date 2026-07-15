package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poorna-9/goshow/internal/services"
)

type ScreenHandler struct {
	Service *services.ScreenService
}

func NewScreenHandler(service *services.ScreenService) *ScreenHandler {
	return &ScreenHandler{Service: service}
}

func (h *ScreenHandler) CreateScreen(c *gin.Context) {
	var req services.CreateScreenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.Service.CreateScreen(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "screen created successfully",
	})
}

func (h *ScreenHandler) GetScreensByTheatre(c *gin.Context) {
	theatreID := c.Param("theatre_id")

	screens, err := h.Service.GetScreensByTheatre(theatreID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, screens)
}

func (h *ScreenHandler) GetScreenInfo(c *gin.Context) {
	id := c.Param("id")
	screen, err := h.Service.GetScreenByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "screen not found"})
		return
	}
	c.JSON(http.StatusOK, screen)
}

func (h *ScreenHandler) DeactivateScreen(c *gin.Context) {
	id := c.Param("id")

	if err := h.Service.DeactivateScreen(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "screen deactivated"})
}
