package models

import (
    "gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
    ID        uint   `gorm:"primaryKey;autoIncrement;column:user_id"`
    FirstName     string `gorm:"size:100"`
    LastName      string `gorm:"size:100"`
    Email         string `gorm:"size:255;unique;not null"`
    Password      string `gorm:"size:255;not null"`
    Phone         string `gorm:"size:20"`
    Role          string `gorm:"size:50;not null"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
    BlockedStatus bool   `gorm:"default:false"`
    InactiveStatus bool  `gorm:"default:false"`
    SuspendedAt    *time.Time
}

type AuthInput struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
}


type OTP struct {
    ID        uint      `gorm:"primaryKey;autoIncrement;column:id"`
    UserID    uint      `gorm:"not null;index"`
    OTP       string    `gorm:"size:10;not null"`
    Expiry    time.Time `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Used      bool      `gorm:"default:false"`
}

type UserSession struct {
	SessionID uint      `gorm:"primaryKey;autoIncrement;session_id"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"not null"`
	Role      string    `gorm:"not null"` 
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time `gorm:"not null"`
}

