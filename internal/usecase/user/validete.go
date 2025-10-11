package user

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

func Validate(req *user.RegisterRequest) map[string]string {
	errors := make(map[string]string)

	// Валидация email
	if err := validateEmail(req.Email); err != nil {
		errors["email"] = err.Error()
	}

	// Валидация пароля
	if err := validatePassword(req.Password); err != nil {
		errors["password"] = err.Error()
	}

	// Проверка совпадения паролей
	if req.Password != req.ConfirmPassword {
		errors["confirm_password"] = "passwords do not match"
	}

	return errors
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)

	// Базовая проверка формата
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	// Проверка длины
	if len(email) > 254 {
		return fmt.Errorf("email too long")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
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
		return fmt.Errorf("password must contain at least 3 of: uppercase, lowercase, numbers, special characters")
	}

	return nil
}
