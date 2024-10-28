package usecase

import (
	"log"
	"time"

	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	GetAllUsers() ([]models.User, error)
	BlockUser(userID uint) error
	UnblockUser(userID uint) error
	GetUserByID(userID uint) (*models.User, error)
	UnsuspendInactiveUsers() error
    UpdateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	CreateSession(userID uint, token string, expiresAt time.Time) error
    DeleteUserSession(sessionID int) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetAllUsers() ([]models.User, error) {
	return u.userRepo.GetAllUsers()
}

func (u *userUsecase) BlockUser(userID uint) error {
	return u.userRepo.UpdateUserBlockStatus(userID, true)
}

func (u *userUsecase) UnblockUser(userID uint) error {
	return u.userRepo.UpdateUserBlockStatus(userID, false)
}

func (u *userUsecase) GetUserByEmail(email string) (*models.User, error) {
    return u.userRepo.GetUserByEmail(email)
}

func VerifyPassword(inputPassword, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func (u *userUsecase) GetUserByID(userID uint) (*models.User, error) {
    return u.userRepo.GetUserByID(userID)
}

func (u *userUsecase) UnsuspendInactiveUsers() error {
    users, err := u.userRepo.GetInactiveSuspendedUsers()
    if err != nil {
        log.Printf("Error fetching suspended users: %v\n", err)
        return err
    }
    successCount := 0
failureCount := 0
for _, user := range users {
    if user.SuspendedAt != nil && time.Since(*user.SuspendedAt) >= 5*24*time.Hour {
        user.InactiveStatus = false
        user.SuspendedAt = nil
        if err := u.userRepo.UpdateUser(&user); err != nil {
            log.Printf("Error updating user (ID: %d): %v\n", user.ID, err)
            failureCount++
        } else {
            log.Printf("User (ID: %d) has been unsuspended.\n", user.ID)
            successCount++
            }
        }
    }
    log.Printf("Unsuspend task completed: %d success, %d failures\n", successCount, failureCount)
    return nil
}

func (u *userUsecase) UpdateUser(user *models.User) error {
    return u.userRepo.UpdateUser(user)
}

func (u *userUsecase) CreateSession(userID uint, token string, expiresAt time.Time) error {
    session := models.UserSession{
        UserID:    userID,
        Token:     token,
        ExpiresAt: expiresAt,
    }
    return u.userRepo.CreateSession(&session)
}

func (u *userUsecase) DeleteUserSession(sessionID int) error {
    return u.userRepo.DeleteSession(uint(sessionID))
}