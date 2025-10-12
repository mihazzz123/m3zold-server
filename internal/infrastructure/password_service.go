package infrastructure

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// PasswordService сервис для работы с паролями
type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (s *PasswordService) ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	// Базовая проверка формата
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return user.ErrInvalidEmail
	}

	// Проверка длины
	if len(email) > 254 {
		return fmt.Errorf("email too long")
	}

	return nil
}

func (s *PasswordService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return user.ErrWeakPassword
	}

	if len(password) > 72 { // bcrypt limitation
		return fmt.Errorf("password too long")
	}

	// Проверка сложности
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	complexity := 0
	if hasUpper {
		complexity++
	}
	if hasLower {
		complexity++
	}
	if hasNumber {
		complexity++
	}
	if hasSpecial {
		complexity++
	}

	if complexity < 3 {
		return user.ErrWeakPassword
	}

	return nil
}

func (s *PasswordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func (s *PasswordService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *PasswordService) GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *PasswordService) GenerateVerificationToken() (string, error) {
	return s.GenerateSecureToken(32)
}
