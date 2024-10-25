package usecase

import (
	"github.com/Prototype-1/xtrace/internal/repository"
	"github.com/Prototype-1/xtrace/internal/models"
)

type UserFavoritesUsecase interface {
    AddFavoriteRoute(userID uint, routeID uint) error
    GetUserFavoriteRoutes(userID uint) ([]models.Route, error)
    RemoveFavoriteRoute(userID uint, routeID uint) error
}

type userFavoritesUsecase struct {
    repo repository.UserFavoritesRepository
}

func NewUserFavoritesUsecase(repo repository.UserFavoritesRepository) UserFavoritesUsecase {
    return &userFavoritesUsecase{repo: repo}
}

func (u *userFavoritesUsecase) AddFavoriteRoute(userID uint, routeID uint) error {
    return u.repo.AddFavoriteRoute(userID, routeID)
}

func (u *userFavoritesUsecase) GetUserFavoriteRoutes(userID uint) ([]models.Route, error) {
    return u.repo.GetUserFavoriteRoutes(userID)
}

func (u *userFavoritesUsecase) RemoveFavoriteRoute(userID uint, routeID uint) error {
    return u.repo.RemoveFavoriteRoute(userID, routeID)
}
