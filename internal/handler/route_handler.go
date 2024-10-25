package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/usecase"
    "net/http"
    "strconv"
)

type RouteHandler struct {
    RouteUsecase *usecase.RouteUsecase
}

func NewRouteHandler(usecase *usecase.RouteUsecase) *RouteHandler {
    return &RouteHandler{RouteUsecase: usecase}
}

func (h *RouteHandler) AddRoute(c *gin.Context) {
    var route models.Route
    if err := c.ShouldBindJSON(&route); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.RouteUsecase.AddRoute(route); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Route added successfully"})
}

func (h *RouteHandler) UpdateRoute(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var route models.Route
    if err := c.ShouldBindJSON(&route); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    route.RouteID = id
    if err := h.RouteUsecase.UpdateRoute(id, route); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update route"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Route updated successfully"})
}

func (h *RouteHandler) DeleteRoute(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.RouteUsecase.DeleteRoute(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete route"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Route deleted successfully"})
}

func (h *RouteHandler) GetAllRoutes(c *gin.Context) {
    routes, err := h.RouteUsecase.GetAllRoutes()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch routes"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"routes": routes})
}

func (h *RouteHandler) GetAllRoutesUser(c *gin.Context) {
    categoryName := c.Query("category")

    if categoryName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Category name is required"})
        return
    }
    routes, err := h.RouteUsecase.GetAllRoutesByCategory(categoryName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch routes"})
        return
    }

    if categoryName == "Rental" {
        c.JSON(http.StatusOK, gin.H{"message": "Service will be available soon"})
        return
    }

    type UserRouteResponse struct {
        RouteID   int    `json:"route_id"`
        RouteName string `json:"route_name"`
    }

    userRoutes := make([]UserRouteResponse, len(routes))
    for i, route := range routes {
        userRoutes[i] = UserRouteResponse{
            RouteID:   route.RouteID,
            RouteName: route.RouteName,
        }
    }

    c.JSON(http.StatusOK, userRoutes)
}

