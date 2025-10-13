package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/mihazzz123/m3zold-server/internal/domain/services"
)

// IDService реализация IDGenerator
type IDService struct{}

// NewIDService создает новый IDService
func NewIDService() services.IDService {
	return &IDService{}
}

// Generate генерирует UUID
func (s *IDService) Generate() string {
	return uuid.New().String()
}

// GenerateSecureToken генерирует безопасный токен
func (s *IDService) GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
