package models

import (
	"time"
)

type Wallet struct {
    WalletID   uint      `gorm:"primaryKey;autoIncrement" json:"wallet_id"`
    UserID     uint      `gorm:"not null" json:"user_id"`
    Balance    float64   `gorm:"not null;default:0" json:"balance"`
    CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type WalletTransaction struct {
    TransactionID uint      `gorm:"primaryKey;autoIncrement" json:"transaction_id"`
    WalletID      uint      `gorm:"not null" json:"wallet_id"`
    AdminID       *uint      `gorm:"not null" json:"admin_id"` 
    Amount        float64   `gorm:"not null" json:"amount"`
    TransactionType          string    `gorm:"not null" json:"type"`  
    Description   string    `gorm:"size:255" json:"description"`
    CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}