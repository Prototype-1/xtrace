package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/usecase"
)

type StopHandler struct {
	StopUsecase *usecase.StopUsecase
}

func NewStopHandler(usecase *usecase.StopUsecase) *StopHandler {
	return &StopHandler{StopUsecase: usecase}
}

func (h *StopHandler) AddStop(c *gin.Context) {
	var stop models.Stop
	if err := c.ShouldBindJSON(&stop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.StopUsecase.AddStop(stop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add stop"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Stop added successfully"})
}

func (h *StopHandler) UpdateStop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var stop models.Stop
	if err := c.ShouldBindJSON(&stop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stop.StopID = id
	if err := h.StopUsecase.UpdateStop(id, stop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stop"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stop updated successfully"})
}

func (h *StopHandler) DeleteStop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.StopUsecase.DeleteStop(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete stop"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stop deleted successfully"})
}

func (h *StopHandler) GetAllStops(c *gin.Context) {
	stops, err := h.StopUsecase.GetAllStops()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stops"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stops": stops})
}
