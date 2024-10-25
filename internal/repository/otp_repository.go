package repository

import (
    "github.com/Prototype-1/xtrace/config"
    "github.com/Prototype-1/xtrace/internal/models"
    "time"
)

// CreateOTP creates a new OTP record in the database
func CreateOTP(otp models.OTP) error {
    return config.DB.Create(&otp).Error
}

// GetOTPByUserID retrieves the OTP record for a specific user by their user ID
func GetOTPByUserID(userID uint) (*models.OTP, error) {
    var otp models.OTP
    err := config.DB.Where("user_id = ? AND used = false AND expiry > ?", userID, time.Now()).First(&otp).Error
    return &otp, err
}

// **NEW** GetOTPByEmail retrieves the OTP record using the user's email
func GetOTPByEmail(email string) (*models.OTP, error) {
    var otp models.OTP
    err := config.DB.Joins("JOIN users ON users.id = otps.user_id").
        Where("users.email = ? AND otps.used = false AND otps.expiry > ?", email, time.Now()).
        First(&otp).Error
    return &otp, err
}

// MarkOTPAsUsed marks an OTP record as used
func MarkOTPAsUsed(otpID uint) error {
    return config.DB.Model(&models.OTP{}).Where("id = ?", otpID).Update("used", true).Error
}

// DeleteExpiredOTPs removes expired OTPs from the database
func DeleteExpiredOTPs() error {
    return config.DB.Where("expiry < ?", time.Now()).Delete(&models.OTP{}).Error
}
