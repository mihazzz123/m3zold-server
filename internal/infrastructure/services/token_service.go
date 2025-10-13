package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/mihazzz123/m3zold-server/internal/domain/services"
)

type TokenService struct{}

func NewTokenService() services.TokenService {
	return &TokenService{}
}

func (s *TokenService) GenerateToken() (string, error) {
	return s.GenerateTokenWithLength(32)
}

func (s *TokenService) GenerateTokenWithLength(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
