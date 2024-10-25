package repository

import (
	"github.com/Prototype-1/xtrace/internal/models"
	"gorm.io/gorm"
)

type StopRepository interface {
	AddStop(stop models.Stop) error
	UpdateStop(stop models.Stop) error
	DeleteStop(id int) error
	GetAllStops() ([]models.Stop, error)
}

type StopRepositoryImpl struct {
	DB *gorm.DB
}

func NewStopRepository(db *gorm.DB) *StopRepositoryImpl {
	return &StopRepositoryImpl{DB: db}
}

func (r *StopRepositoryImpl) AddStop(stop models.Stop) error {
	return r.DB.Create(&stop).Error
}

func (r *StopRepositoryImpl) UpdateStop(stop models.Stop) error {
    var existingStop models.Stop

    if err := r.DB.First(&existingStop, stop.StopID).Error; err != nil {
        return err
    }
    return r.DB.Model(&existingStop).Updates(map[string]interface{}{
        "stop_name":   stop.StopName,
        "latitude":    stop.Latitude,
        "longitude":   stop.Longitude,
        "category_id": stop.CategoryID,
        "updated_at":  gorm.Expr("NOW()"), 
    }).Error
}

func (r *StopRepositoryImpl) DeleteStop(id int) error {
	var stop models.Stop
	if err := r.DB.First(&stop, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&stop).Error
}

func (r *StopRepositoryImpl) GetAllStops() ([]models.Stop, error) {
	var stops []models.Stop
	err := r.DB.Find(&stops).Error
	return stops, err
}
