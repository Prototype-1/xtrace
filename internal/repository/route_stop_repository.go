package repository

import (
	"github.com/Prototype-1/xtrace/internal/models"
	"gorm.io/gorm"
	"fmt"
)

type RouteStopRepository interface {
	AddRouteStop(routeStop models.RouteStop) error
	UpdateRouteStop(routeStop models.RouteStop) error
	DeleteRouteStop(id int) error
	GetAllRouteStops() ([]models.RouteStop, error)
	GetOrderedStopsByRouteID(routeID uint) ([]models.RouteStop, error)
	GetStopByID(int) (*models.Stop, string, error)
	GetStopsByRouteID(routeID uint) ([]models.Stop, error)
}

type RouteStopRepositoryImpl struct {
	DB *gorm.DB
}

func NewRouteStopRepository(db *gorm.DB) *RouteStopRepositoryImpl {
	return &RouteStopRepositoryImpl{DB: db}
}

func (r *RouteStopRepositoryImpl) AddRouteStop(routeStop models.RouteStop) error {
	return r.DB.Create(&routeStop).Error
}

func (r *RouteStopRepositoryImpl) UpdateRouteStop(routeStop models.RouteStop) error {
    return r.DB.Model(&routeStop).Select("stop_id", "stop_sequence", "updated_at").Updates(routeStop).Error
}


func (r *RouteStopRepositoryImpl) DeleteRouteStop(id int) error {
	var routeStop models.RouteStop
	if err := r.DB.First(&routeStop, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&routeStop).Error
}

func (r *RouteStopRepositoryImpl) GetAllRouteStops() ([]models.RouteStop, error) {
	var routeStops []models.RouteStop
	err := r.DB.Find(&routeStops).Error
	return routeStops, err
}

func (r *RouteStopRepositoryImpl) GetOrderedStopsByRouteID(routeID uint) ([]models.RouteStop, error) {
    var routeStops []models.RouteStop
    query := `
        SELECT * 
        FROM route_stops
        WHERE route_id = ?
        ORDER BY stop_sequence;
    `
    if err := r.DB.Raw(query, routeID).Scan(&routeStops).Error; err != nil {
        return nil, fmt.Errorf("error retrieving ordered route stops: %w", err)
    }

    return routeStops, nil
}

func (r *RouteStopRepositoryImpl) GetStopByID(id int) (*models.Stop, string, error) {
    var stop models.Stop
    var categoryName string

    // Retrieve the stop
    err := r.DB.Where("stop_id = ?", id).First(&stop).Error
    if err != nil {
        return nil, "", err
    }

    var category models.Category
    err = r.DB.Where("category_id = ?", stop.CategoryID).First(&category).Error
    if err != nil {
        return &stop, "", nil
    }

    categoryName = category.CategoryName

    return &stop, categoryName, nil
}


func (r *RouteStopRepositoryImpl) GetStopsByRouteID(routeID uint) ([]models.Stop, error) {
    var stops []models.Stop
    query := `
        SELECT s.*
        FROM stops s
        JOIN route_stops rs ON s.stop_id = rs.stop_id
        WHERE rs.route_id = ?
    `
    if err := r.DB.Raw(query, routeID).Scan(&stops).Error; err != nil {
        return nil, fmt.Errorf("error retrieving stops for route ID %d: %w", routeID, err)
    }
    return stops, nil
}




