package models

import "time"

type Invoice struct {
    InvoiceID      uint      `gorm:"primaryKey"`
    UserID         uint      `gorm:"not null"`
    PaymentID      uint      `gorm:"not null"`
    InvoiceDate    time.Time `gorm:""`
    OriginalAmount float64   `gorm:"not null"`
    DiscountAmount float64   `gorm:"default:0"`
    Amount         float64   `gorm:"not null"`
    Status         string    `gorm:"size:50"`
    PaymentType    string    `gorm:"not null"`
    CreatedAt      time.Time `gorm:"autoCreateTime"`
    UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

