package models

import "time"

type Category struct {
    CategoryID        int       `json:"category_id" gorm:"column:category_id;primary_key"`
    CategoryName     string    `json:"category_name" validate:"required"`
    ImageURL     string    `json:"image_url"`
    IsDeleted bool      `json:"is_deleted"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}