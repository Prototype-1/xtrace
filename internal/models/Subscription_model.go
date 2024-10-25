package models

import (
    "time"
)

type Subscription struct {
    SubscriptionID uint      `gorm:"primaryKey;autoIncrement" json:"subscription_id"`
    UserID         uint      `json:"user_id"` 
    PlanID         uint      `json:"plan_id"`
    StartDate      time.Time `json:"start_date"`
    EndDate        time.Time `json:"end_date"`  
    ServiceType    string    `gorm:"type:varchar(50)" json:"service_type"` 
     CardType       string    `json:"card_type"`
     NolCardID      uint      `json:"nol_card_id"`
	Price          float64       `json:"price"`
	DurationDays       int       `json:"duration_days"`
    CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type SubscriptionPlan struct {
    PlanID       uint      `gorm:"primaryKey" json:"plan_id"`
    PlanName     string    `json:"plan_name"`
    Price        float64   `json:"price"`
    DurationDays int       `json:"duration_days"`
    CardType      string    `json:"card_type"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
