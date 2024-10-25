package handler

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/Prototype-1/xtrace/config"
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/pkg/utils"
	"github.com/Prototype-1/xtrace/internal/repository"
	"strconv"
)

func AdminSignUp(c *gin.Context) {
	var input models.AuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
        return
    }

	existingUser, err := repository.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	admin := models.User{
		Email:    input.Email,
		Password: hashedPassword,
		Role:     "admin",
	}

	err = repository.CreateUser(&admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin account created successfully"})
}

func AdminLogin(c *gin.Context) {
	var input models.AuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := repository.GetUserByEmail(input.Email)
	if err != nil || user == nil || user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Create JWT tokens
	tokenDetails, err := utils.CreateToken(int(user.ID), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	session := models.UserSession{
		UserID:    user.ID,
		Token:     tokenDetails.AccessToken,
		Role:      user.Role,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), 
	}


	sessionRepo := repository.NewUserRepository(config.DB) 
	err = sessionRepo.CreateSession(&session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "You are successfully logged in as Admin"})
	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	})
}

func SuspendUser(c *gin.Context) {
	var input struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repository.NewUserRepository(config.DB)

	user, err := userRepo.GetUserByID(input.UserID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	now := time.Now()
	user.InactiveStatus = true
	user.SuspendedAt = &now

	if err := userRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User suspended successfully", "suspended_until": now.AddDate(0, 0, 5)})
}

func AdminLogout(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not an Admin for this."})
		return
	}
	sessionIDStr := c.Query("session_id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID, You are not Logged In."})
		return
	}
	userRepo := repository.NewUserRepository(config.DB)
	if err := userRepo.DeleteSession(uint(sessionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}


