package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/Prototype-1/xtrace/config"
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/repository"
	"github.com/Prototype-1/xtrace/internal/usecase"
	"github.com/Prototype-1/xtrace/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	googleOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func UserSignUp(c *gin.Context) {
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

    user := models.User{
        FirstName: input.FirstName,
        LastName:  input.LastName,
        Email:     input.Email,
        Password:  hashedPassword,
        Role:      "user",
    }

    err = repository.CreateUser(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    otpValue := utils.GenerateOTP()
    otp := models.OTP{
        UserID: user.ID,
        OTP:    otpValue,
        Expiry: time.Now().Add(5 * time.Minute),
    }
    err = repository.CreateOTP(otp)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OTP"})
        return
    }

    err = utils.SendEmail(user.Email, "Your OTP Code", fmt.Sprintf("Your OTP code is: %s", otpValue))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "OTP is shared through your email. Verify to complete sign-up."})
}

func VerifyOTP(c *gin.Context) {
    var input struct {
        Email string `json:"email" binding:"required"`
        OTP   string `json:"otp" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Retrieve user by email
    user, err := repository.GetUserByEmail(input.Email)
    if err != nil || user == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
        return
    }

    // Retrieve the OTP
    otp, err := repository.GetOTPByUserID(user.ID)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP required or expired"})
        return
    }

    if otp.OTP != input.OTP || otp.Used || time.Now().After(otp.Expiry) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
        return
    }
    err = repository.MarkOTPAsUsed(otp.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully. Your account is now active."})
}

func UserLogin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repository.NewUserRepository(config.DB)
	user, err := userRepo.GetUserByEmail(input.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.BlockedStatus {
		c.JSON(http.StatusForbidden, gin.H{"error": "Your account is blocked due to policy violations."})
		return
	}

	if user.InactiveStatus && user.SuspendedAt != nil {
		suspensionEnd := user.SuspendedAt.AddDate(0, 0, 5)
		if time.Now().Before(suspensionEnd) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Your account is temporarily suspended. Please contact support."})
			return
		} else {
			user.InactiveStatus = false
			user.SuspendedAt = nil
			if err := userRepo.UpdateUser(user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
				return
			}
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
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
		ExpiresAt: time.Now().Add(24 * time.Hour), // Session valid for 24 hours
	}

	if err := userRepo.CreateSession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

    c.JSON(http.StatusOK, gin.H{"message": "Login Successful."})

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	})
}


func ResendOTP(c *gin.Context) {
    var input struct {
        Email string `json:"email" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get user by email
    user, err := repository.GetUserByEmail(input.Email)
    if err != nil || user == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Generate a new OTP
    otpValue := utils.GenerateOTP()
    otp := models.OTP{
        UserID: user.ID,
        OTP:    otpValue,
        Expiry: time.Now().Add(5 * time.Minute),
    }

    // Update or create OTP in the database
    err = repository.UpdateOrCreateOTP(otp)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new OTP"})
        return
    }

    // Send OTP via email
    err = utils.SendEmail(user.Email, "Your OTP Code", fmt.Sprintf("Your new OTP code is: %s", otpValue))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func GoogleLogin(c *gin.Context) {
    url := config.GoogleOAuthConfig.AuthCodeURL("random-state")
	fmt.Println("Authorization URL:", url)
    c.JSON(http.StatusOK, gin.H{"authorization_url": url})
}

// GoogleCallback handles the callback from Google
func GoogleCallback(c *gin.Context) {
    state := c.Query("state")
    if state != "random-state" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "State mismatch"})
        return
    }

    code := c.Query("code")
    token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
    if err != nil {
		log.Println("Error during token exchange:", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to exchange token"})
        return
    }

    // Get user info from Google
	client := config.GoogleOAuthConfig.Client(context.Background(), token)
    
    // Use NewService and option.WithHTTPClient to pass the client
    service, err := googleOAuth2.NewService(context.Background(), option.WithHTTPClient(client))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OAuth service"})
        return
    }

    userInfo, err := service.Userinfo.Get().Do()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
        return
    }
    fmt.Println(userInfo)

    user, err := repository.GetUserByEmail(userInfo.Email)
    if err != nil || user == nil {
        hashedPassword, _ := utils.HashPassword("dummy_password") 
        newUser := models.User{
            FirstName:  userInfo.GivenName,
            LastName:   userInfo.FamilyName,
            Email:      userInfo.Email,
            Password:   hashedPassword,  
            Role:       "user",         
            BlockedStatus: false,       
            InactiveStatus: false,      
        }
        err = repository.CreateUser(&newUser)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new user"})
            return
        }
        user = &newUser
    }

    tokenDetails, err := utils.CreateToken(int(user.ID), user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Return user info and token
    c.JSON(http.StatusOK, gin.H{
        "token":     tokenDetails.AccessToken,
        "email":     user.Email,
        "name":      user.FirstName + " " + user.LastName,
        "role":      user.Role,
    })
}

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch users"})
		return
	}
    totalUsers := len(users)

	c.JSON(http.StatusOK, gin.H{
		"total": totalUsers,
		"users": users,
	})
}

func (h *UserHandler) BlockUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.userUsecase.BlockUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to block user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

// Unblocks a user
func (h *UserHandler) UnblockUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.userUsecase.UnblockUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unblock user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User unblocked successfully"})
}

func (h *UserHandler) SuspendUser(c *gin.Context) {
    userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    user, err := h.userUsecase.GetUserByID(uint(userID))
    if err != nil || user == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    now := time.Now()
    user.InactiveStatus = true
    user.SuspendedAt = &now
    if err := h.userUsecase.UpdateUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend user"})
        return
    }
    suspendedUntil := now.Add(5 * 24 * time.Hour)

    c.JSON(http.StatusOK, gin.H{
        "message":         "User suspended successfully",
        "suspended_until": suspendedUntil,
    })
}

func UserLogout(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Your are not a ordinary user, use your admin log out."})
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

func (h *UserHandler) GetUserStatus(c *gin.Context) {
    userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    user, err := h.userUsecase.GetUserByID(uint(userID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "blocked": user.BlockedStatus,
        "inactive": user.InactiveStatus,
    })
}


