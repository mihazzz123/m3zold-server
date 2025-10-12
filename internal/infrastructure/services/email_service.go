package infrastructure_services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

type EmailService struct {
	emailRegex *regexp.Regexp
}

func NewEmailService() *EmailService {
	// Простая валидация email, можно использовать более сложную
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return &EmailService{emailRegex: emailRegex}
}

func (s *EmailService) Validate(email string) error {
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

func (s *EmailService) Normalize(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
