package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"github.com/Prototype-1/xtrace/internal/models"
)

func CreateToken(userID int, role string) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}
	var err error

	// Access token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["role"] = role
	atClaims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	// Refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}
