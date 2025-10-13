package services

import (
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
