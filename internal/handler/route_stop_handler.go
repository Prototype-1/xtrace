package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/usecase"
	"net/http"
	"strconv"
)

type RouteStopHandler struct {
	RouteStopUsecase *usecase.RouteStopUsecase
}

func NewRouteStopHandler(usecase *usecase.RouteStopUsecase) *RouteStopHandler {
	return &RouteStopHandler{RouteStopUsecase: usecase}
}

func (h *RouteStopHandler) AddRouteStop(c *gin.Context) {
	var routeStop models.RouteStop
	if err := c.ShouldBindJSON(&routeStop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.RouteStopUsecase.AddRouteStop(routeStop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add route stop"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Route stop added successfully"})
}

func (h *RouteStopHandler) UpdateRouteStop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var routeStop models.RouteStop
	if err := c.ShouldBindJSON(&routeStop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	routeStop.RouteStopID = id
	if err := h.RouteStopUsecase.UpdateRouteStop(id, routeStop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update route stop"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Route stop updated successfully"})
}

func (h *RouteStopHandler) DeleteRouteStop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.RouteStopUsecase.DeleteRouteStop(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete route stop"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Route stop deleted successfully"})
}

func (h *RouteStopHandler) GetAllRouteStops(c *gin.Context) {
	routeStops, err := h.RouteStopUsecase.GetAllRouteStops()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch route stops"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"route_stops": routeStops})
}

func (h *RouteStopHandler) GetOrderedStopsByRoute(c *gin.Context) {
    routeIDParam := c.Param("route_id")
    routeID, err := strconv.Atoi(routeIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route id"})
        return
    }

    category := c.Query("category")
    if category != "Metro" && category != "Bus" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Only 'Metro' or 'Bus' categories are allowed"})
        return
    }

    orderedStops, err := h.RouteStopUsecase.GetOrderedStopsByCategory(uint(routeID), category)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"stops": orderedStops})
}


func (h *RouteStopHandler) FindNearestStop(c *gin.Context) {
	var input struct {
		Latitude  float64 `json:"latitude" binding:"required,numeric"`
		Longitude float64 `json:"longitude" binding:"required,numeric"`
		RouteID   int     `json:"route_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nearestStop, err := h.RouteStopUsecase.FindNearestStop(input.Latitude, input.Longitude, input.RouteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := struct {
		StopID   uint   `json:"stop_id"`
		StopName string `json:"stop_name"`
		Latitude float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{
		StopID:   uint(nearestStop.StopID),
		StopName: nearestStop.StopName,
		Latitude: nearestStop.Latitude,
		Longitude: nearestStop.Longitude,
	}

	c.JSON(http.StatusOK, gin.H{"nearest_stop": response})
}

