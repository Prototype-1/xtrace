package repository

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "gorm.io/gorm"
)

type FareRuleRepository interface {
    CreateFareRule(fareRule models.FareRule) error
    UpdateFareRule(fareRule models.FareRule) error
    DeleteFareRule(id int) error
    GetFareRuleByID(id int) (models.FareRule, error)
    GetAllFareRules() ([]models.FareRule, error)
    GetFareRuleByRouteID(routeID int) (models.FareRule, error)
    //GetTravelTimeBetweenStops(fromStopID, toStopID uint) (int, error)
    CreateStopDuration(stopDuration models.StopDuration) error
    UpdateStopDuration(stopDuration models.StopDuration) error
    DeleteStopDuration(id uint) error
    GetAllStopDurations() ([]models.StopDuration, error)
}

type FareRuleRepositoryImpl struct {
    DB *gorm.DB
}

func NewFareRuleRepository(db *gorm.DB) FareRuleRepository {
    return &FareRuleRepositoryImpl{DB: db}
}

func (r *FareRuleRepositoryImpl) CreateFareRule(fareRule models.FareRule) error {
    return r.DB.Create(&fareRule).Error
}

func (r *FareRuleRepositoryImpl) UpdateFareRule(fareRule models.FareRule) error {
    return r.DB.Model(&fareRule).Select("base_fare", "fare_per_km", "fare_per_stop", "base_km", "base_stops", "updated_at").Updates(fareRule).Error
}

func (r *FareRuleRepositoryImpl) DeleteFareRule(id int) error {
    return r.DB.Delete(&models.FareRule{}, id).Error
}

func (r *FareRuleRepositoryImpl) GetFareRuleByID(id int) (models.FareRule, error) {
    var fareRule models.FareRule
    err := r.DB.First(&fareRule, id).Error
    return fareRule, err
}

func (r *FareRuleRepositoryImpl) GetAllFareRules() ([]models.FareRule, error) {
    var fareRules []models.FareRule
    err := r.DB.Find(&fareRules).Error
    return fareRules, err
}

func (r *FareRuleRepositoryImpl) GetFareRuleByRouteID(routeID int) (models.FareRule, error) {
    var fareRule models.FareRule
    err := r.DB.Where("route_id = ?", routeID).First(&fareRule).Error
    return fareRule, err
}

func (r *FareRuleRepositoryImpl) CreateStopDuration(stopDuration models.StopDuration) error {
    return r.DB.Create(&stopDuration).Error
}

func (r *FareRuleRepositoryImpl) UpdateStopDuration(stopDuration models.StopDuration) error {
    return r.DB.Model(&stopDuration).Updates(stopDuration).Error
}

func (r *FareRuleRepositoryImpl) DeleteStopDuration(id uint) error {
    return r.DB.Delete(&models.StopDuration{}, id).Error
}

func (r *FareRuleRepositoryImpl) GetAllStopDurations() ([]models.StopDuration, error) {
    var stopDurations []models.StopDuration
    err := r.DB.Find(&stopDurations).Error
    return stopDurations, err
}

