package usecase

import (
	"github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/repository"
)

type RouteUsecase struct {
    RouteRepo repository.RouteRepository
}

func NewRouteUsecase(repo repository.RouteRepository) *RouteUsecase {
    return &RouteUsecase{RouteRepo: repo}
}

func (u *RouteUsecase) AddRoute(route models.Route) error {
    return u.RouteRepo.AddRoute(route)
}

func (u *RouteUsecase) UpdateRoute(id int, route models.Route) error {
    return u.RouteRepo.UpdateRoute(route)
}

func (u *RouteUsecase) DeleteRoute(id int) error {
    return u.RouteRepo.DeleteRoute(id)
}

func (u *RouteUsecase) GetAllRoutes() ([]models.Route, error) {
    return u.RouteRepo.GetAllRoutes()
}

func (u *RouteUsecase) GetAllRoutesByCategory(categoryName string) ([]models.Route, error) {
    return u.RouteRepo.GetAllRoutesByCategory(categoryName)
}