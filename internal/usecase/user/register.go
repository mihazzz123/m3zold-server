package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mihazzz123/m3zold-server/internal/constants"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	UserName        string `json:"user_name"`
}
type RegisterUseCase struct {
	Repo user.Repository
}

func NewRegisterUseCase(repo user.Repository) *RegisterUseCase {
	return &RegisterUseCase{Repo: repo}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterRequest) (error, &user.User) {
	// Валидация входных данных
	if errors := input.Validate(); len(errors) > 0 {
		return fmt.Errorf("validation failed: %v", errors)
	}

	exists, err := uc.Repo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if exists {
		return constants.ErrEmailTaken
	}

	// Хеширование пароля
	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("failed to process password: %w", err)
	}

	// Создание пользователя
	user := &user.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     false, // Требуется верификация
	}

	if err = uc.Repo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Генерация токена верификации
	verificationToken, err := generateVerificationToken()
	if err != nil {
		// Логируем ошибку, но не прерываем процесс
		// Пользователь может запросить повторную отправку
		fmt.Printf("Failed to generate verification token: %v\n", err)
	} else {
		// Отправка email с токеном (в реальном приложении)
		go s.sendVerificationEmail(user.Email, verificationToken)
	}

	// Не возвращаем хеш пароля
	user.PasswordHash = ""

	return user, nil
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

func (req *RegisterRequest) Validate() map[string]string {
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

func hashPassword(password string) (string, error) {
	// Генерация соли и хеширование
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func generateVerificationToken() (string, error) {
	return generateSecureToken(32)
}
