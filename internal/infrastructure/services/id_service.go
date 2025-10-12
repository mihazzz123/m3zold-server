package infrastructure_services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

type IDService struct{}

func NewIDService() *IDService {
	return &IDService{}
}

func (s *IDService) Generate() string {
	return uuid.New().String()
}

func (s *IDService) GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
