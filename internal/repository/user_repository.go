package repository

import (
	"errors"
	"log"
	"time"
	"strings"
	"gorm.io/gorm"
	"github.com/Prototype-1/xtrace/config"
	"github.com/Prototype-1/xtrace/internal/models"
)

func CreateUser(user *models.User) error {
	if config.DB == nil {
		return errors.New("database connection is not initialized")
	}

	result := config.DB.Create(user)
	if result.Error != nil {
		log.Printf("Error creating user: %v\n", result.Error)
		return result.Error
	}
	log.Printf("User created: %+v\n", user)
	return nil
}

// Retrieves a user by their email
func GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    result := config.DB.Where("email = ?", email).First(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil 
        }
        return nil, result.Error
    }
    return &user, nil
}


func UpdateOrCreateOTP(otp models.OTP) error {
	var existingOTP models.OTP
	err := config.DB.Where("user_id = ?", otp.UserID).First(&existingOTP).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingOTP.ID != 0 {
		existingOTP.OTP = otp.OTP
		existingOTP.Expiry = otp.Expiry
		existingOTP.Used = false
		return config.DB.Save(&existingOTP).Error
	}
	return config.DB.Create(&otp).Error
}

type UserRepository interface {
	GetAllUsers() ([]models.User, error)
	UpdateUserBlockStatus(userID uint, blockedStatus bool) error
	UpdateUserInactiveStatus(userID uint, inactiveStatus bool, suspendedAt *time.Time) error
	GetUserByID(userID uint) (*models.User, error)
	GetInactiveSuspendedUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	CreateSession(session *models.UserSession) error
	DeleteSession(sessionID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// UpdateUserBlockStatus updates the blocked status of a user
func (r *userRepository) UpdateUserBlockStatus(userID uint, blockedStatus bool) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}
	user.BlockedStatus = blockedStatus
	return r.db.Save(&user).Error
}

func (r *userRepository) UpdateUserInactiveStatus(userID uint, inactiveStatus bool, suspendedAt *time.Time) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}
	user.InactiveStatus = inactiveStatus
	user.SuspendedAt = suspendedAt
	return r.db.Save(&user).Error
}


func (r *userRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetInactiveSuspendedUsers() ([]models.User, error) {
    var users []models.User
    err := r.db.Where("inactive_status = ? AND suspended_at IS NOT NULL", true).Find(&users).Error
    return users, err
}


func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil 
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) CreateSession(session *models.UserSession) error {
	return r.db.Create(session).Error
}

func (r *userRepository) DeleteSession(sessionID uint) error {
	result := r.db.Delete(&models.UserSession{}, sessionID)
	if result.RowsAffected == 0 {
		return errors.New("no session found with the provided session ID")
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}



