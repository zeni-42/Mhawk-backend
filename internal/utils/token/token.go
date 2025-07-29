package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeni-42/Mhawk/internal/models"
)

func GetAccessToken(u models.User) string {
	claims := jwt.MapClaims{
		"sub": u.Id.String(),
		"iss": "Mhawk",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessSecret := os.Getenv("ACCESS_SECRET")

	signedToken, err := token.SignedString([]byte(accessSecret))
	if err != nil {
		return ""
	}

	return signedToken
}

func GetRefreshToken(u models.User) string {
	refreshData := map[string]interface{}{
		"userid": u.Id,
		"fullName": u.Fullname,
		"isPro": u.IsPro,
	}
	
	claims := jwt.MapClaims{
		"sub": refreshData,
		"iss": "Mhawk",
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshSecret := os.Getenv("REFRESH_SECRET")

	signedToken, err := token.SignedString([]byte(refreshSecret))

	if err != nil {
		return ""
	}

	return signedToken
}