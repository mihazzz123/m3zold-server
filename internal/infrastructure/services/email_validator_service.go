package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mihazzz123/m3zold-server/internal/domain/services"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

// EmailValidatorService реализация EmailValidator
type EmailValidatorService struct {
	emailRegex *regexp.Regexp
}

// NewEmailValidatorService создает новый EmailService
func NewEmailValidatorService() services.EmailValidatorService {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return &EmailValidatorService{emailRegex: emailRegex}
}

// Validate валидирует email
func (s *EmailValidatorService) Validate(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return user.ErrInvalidEmail
	}

	if len(email) > 254 {
		return fmt.Errorf("email too long")
	}

	if !s.emailRegex.MatchString(email) {
		return user.ErrInvalidEmail
	}

	return nil
}

// Normalize нормализует email
func (s *EmailValidatorService) Normalize(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
