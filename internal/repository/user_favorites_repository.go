package repository

import (
	 "gorm.io/gorm"
    "github.com/Prototype-1/xtrace/internal/models"
)

type UserFavoritesRepository interface {
    AddFavoriteRoute(userID uint, routeID uint) error
    GetUserFavoriteRoutes(userID uint) ([]models.Route, error)
    RemoveFavoriteRoute(userID uint, routeID uint) error
}

type userFavoritesRepository struct {
    DB *gorm.DB
}

func NewUserFavoritesRepository(db *gorm.DB) UserFavoritesRepository {
    return &userFavoritesRepository{DB: db}
}

func (r *userFavoritesRepository) AddFavoriteRoute(userID uint, routeID uint) error {
    favorite := models.UserFavorite{UserID: userID, RouteID: routeID}
    return r.DB.Create(&favorite).Error
}

func (r *userFavoritesRepository) GetUserFavoriteRoutes(userID uint) ([]models.Route, error) {
    var routes []models.Route
    err := r.DB.Table("routes").Select("*").
        Joins("JOIN user_favorites ON routes.route_id = user_favorites.route_id").
        Where("user_favorites.user_id = ?", userID).Find(&routes).Error
    return routes, err
}

func (r *userFavoritesRepository) RemoveFavoriteRoute(userID uint, routeID uint) error {
    return r.DB.Where("user_id = ? AND route_id = ?", userID, routeID).Delete(&models.UserFavorite{}).Error
}
