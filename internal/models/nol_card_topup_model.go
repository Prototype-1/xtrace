package models

import "time"

type NolCard struct {
    NolCardID  int       `json:"nol_card_id" gorm:"primary_key;primaryKey;autoIncrement"`
    UserID     int       `json:"user_id"`
    CardNumber string    `json:"card_number"`
    Balance    float64   `json:"balance"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    CardType   string    `json:"card_type"` 
}

type NolCardTopup struct {
    TopupID    int       `json:"topup_id" gorm:"primary_key;primaryKey;autoIncrement"`
    NolCardID  int       `json:"nol_card_id"`
    Amount     float64   `json:"amount"`
    TopupDate  time.Time `json:"topup_date"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
