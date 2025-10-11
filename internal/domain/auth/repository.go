package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mihazzz123/m3zold-server/internal/config"
)

func GenerateToken(cfg *config.Config, userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(cfg.Auth.JWTSecret)
}
