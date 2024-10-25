package models

import "time"

type Coupon struct {
    CouponID        int       `gorm:"primary_key;auto_increment" json:"coupon_id"`
    Code            string    `json:"code"`
    DiscountAmount  float64   `json:"discount_amount"`
    DiscountType    string    `json:"discount_type"`
    StartDate       time.Time `json:"start_date"`
    EndDate         time.Time `json:"end_date"`
    PaymentType   string        `json:"payment_type"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

