package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/Prototype-1/xtrace/config"
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, "Please Login to your account")
			log.Println("Access token not provided")
			c.Abort()
			return
		}

		tokenString = strings.Split(tokenString, "Bearer ")[1]
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, "Invalid access token")
			c.Abort()
			return
		}

		var session models.UserSession
		userID := uint((*claims)["user_id"].(float64))

		if err := config.DB.Where("user_id = ? AND token = ?", userID, tokenString).First(&session).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, "You may have to Login again.")
				log.Println("Session not found for the token:", tokenString)
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, "Database error")
			c.Abort()
			return
		}
		c.Set("user_id", (*claims)["user_id"])
		c.Set("role", (*claims)["role"])
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, "Seems like you don't have the Admin Rights!!!")
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "user" {
			c.JSON(http.StatusForbidden, "Seems like you're are not an ordinary user, please login as the Admin!!!")
			c.Abort()
			return
		}
		c.Next()
	}
}

