package models

import "time"

type Booking struct {
    BookingID     uint      `gorm:"primaryKey;autoIncrement" json:"booking_id"`
    UserID        uint      `json:"user_id"`
    RouteID       uint      `json:"route_id"`
    PaymentID     *uint     `json:"payment_id"` 
    ServiceType   string    `json:"service_type"`
    CardType      string    `json:"card_type"`
    BookingAmount float64   `json:"booking_amount"`
    Status        string    `json:"status"`
    BookingDate   time.Time `json:"booking_date" gorm:"default:CURRENT_TIMESTAMP"`
    CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
    UpdatedAt     time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

