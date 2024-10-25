package models

import (
    "time"
)

type RazorpayPayment struct {
    PaymentID      uint    `gorm:"primaryKey;autoIncrement" json:"payment_id"`
    UserID         uint      `json:"user_id"`
    RazorpayID     string    `json:"razorpay_id"` 
    OrderID        string    `json:"order_id"`    
    Amount         float64   `json:"amount"`      
    Currency       string    `json:"currency"`    
    Status         string    `json:"status"`     
    Method         string    `json:"method"`      
    CouponCode     string    `json:"coupon_code"`
    WalletID       *uint      `json:"wallet_id"`     
    NolCardID     *uint      `json:"nol_card_id"`   
    SubscriptionID *uint      `json:"subscription_id"`
    BookingID      *uint      `json:"booking_id"`   
    CouponID       *uint      `json:"coupon_id"`     
    PaymentType    string    `json:"payment_type"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    PaymentDate    time.Time `json:"payment_date"`
}
