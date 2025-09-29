package utils

import (
	"time"

	"github.com/fahmiarz/project-management/config"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

//generatetoke jwt
//generate refresh token
func GenerateToken(userID int64, role, email string, publicID uuid.UUID) (string, error) {
	secret := config.AppConfig.JWTSecret
	duration, _ := time.ParseDuration(config.AppConfig.JWTExpire)
	claims := jwt.MapClaims{
		"user_id" : userID,
		"role" : role,
		"public_id" : publicID,
		"email" : email,
		"exp" : time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}