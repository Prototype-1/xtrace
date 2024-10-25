package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/Prototype-1/xtrace/internal/usecase"
)

type UserFavoritesHandler struct {
    userFavoritesUsecase usecase.UserFavoritesUsecase
}

func NewUserFavoritesHandler(userFavoritesUsecase usecase.UserFavoritesUsecase) *UserFavoritesHandler {
    return &UserFavoritesHandler{userFavoritesUsecase: userFavoritesUsecase}
}

func (h *UserFavoritesHandler) AddFavoriteRoute(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("userID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    routeID, err := strconv.Atoi(c.Param("routeID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
        return
    }

    err = h.userFavoritesUsecase.AddFavoriteRoute(uint(userID), uint(routeID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add favorite"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Route added to favorites"})
}

func (h *UserFavoritesHandler) GetUserFavoriteRoutes(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("userID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    routes, err := h.userFavoritesUsecase.GetUserFavoriteRoutes(uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch favorite routes"})
        return
    }

    c.JSON(http.StatusOK, routes)
}

func (h *UserFavoritesHandler) RemoveFavoriteRoute(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("userID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    routeID, err := strconv.Atoi(c.Param("routeID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
        return
    }

    err = h.userFavoritesUsecase.RemoveFavoriteRoute(uint(userID), uint(routeID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove favorite"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Route removed from favorites"})
}
