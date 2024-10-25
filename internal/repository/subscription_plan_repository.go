package repository

import (
	"fmt"
	"github.com/Prototype-1/xtrace/internal/models"
	"gorm.io/gorm"
)

type SubscriptionPlanRepository interface {
    CreateSubscriptionPlan(plan *models.SubscriptionPlan) error
    GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error)
    UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error
    DeleteSubscriptionPlan(id uint) error
    GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error)
}

type SubscriptionPlanRepositoryImpl struct {
    DB *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) SubscriptionPlanRepository {
    return &SubscriptionPlanRepositoryImpl{DB: db}
}

func (r *SubscriptionPlanRepositoryImpl) CreateSubscriptionPlan(plan *models.SubscriptionPlan) error {
    return r.DB.Create(plan).Error
}

func (r *SubscriptionPlanRepositoryImpl) GetSubscriptionPlanByID(id uint) (*models.SubscriptionPlan, error) {
    var plan models.SubscriptionPlan
    err := r.DB.First(&plan, id).Error
    if err != nil {
        fmt.Printf("Error fetching plan: %v\n", err) 
        return nil, fmt.Errorf("enter a valid plan ID")
    }
    return &plan, nil
}

func (r *SubscriptionPlanRepositoryImpl) UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error {
    // Updates` to make sure only specified fields are updated
    return r.DB.Model(&models.SubscriptionPlan{PlanID: plan.PlanID}).Updates(map[string]interface{}{
        "PlanName":     plan.PlanName,
        "Price":        plan.Price,
        "DurationDays": plan.DurationDays,
    }).Error
}


func (r *SubscriptionPlanRepositoryImpl) DeleteSubscriptionPlan(id uint) error {
    return r.DB.Delete(&models.SubscriptionPlan{}, id).Error
}

func (r *SubscriptionPlanRepositoryImpl) GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error) {
    var plans []models.SubscriptionPlan
    err := r.DB.Find(&plans).Error
    return plans, err
}
