package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService реализация PasswordHasher
type PasswordService struct {
	cost int
}

// NewPasswordService создает новый PasswordService
func NewPasswordService(cost int) *PasswordService {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &PasswordService{cost: cost}
}

// Hash хеширует пароль
func (s *PasswordService) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// Verify проверяет пароль
func (s *PasswordService) Verify(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
