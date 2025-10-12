package user

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUseCase struct {
	Repo user.Repository
}

func NewRegisterUseCase(repo user.Repository) *RegisterUseCase {
	return &RegisterUseCase{Repo: repo}
}

// RegisterRequest DTO для запроса регистрации
type RegisterRequest struct {
	Email           string `json:"email"`
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
}

// RegisterResponse DTO для ответа регистрации
type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Execute выполняет регистрацию пользователя
func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterRequest) (*RegisterResponse, error) {
	// Валидация входных данных
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Нормализация email
	email := strings.ToLower(strings.TrimSpace(input.Email))

	// Проверка существования пользователя
	exists, err := uc.Repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if exists {
		return nil, user.ErrEmailTaken
	}

	// Хеширование пароля
	passwordHash, err := uc.hashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("password processing error: %w", err)
	}

	// Создание пользователя
	newUser := &user.User{
		ID:           uuid.New().String(),
		Email:        email,
		UserName:     input.UserName,
		PasswordHash: passwordHash,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    nil,
	}

	// Сохранение в репозиторий
	if err = uc.Repo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := uuid.Parse(newUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}

	// Возвращаем ответ без чувствительных данных
	return &RegisterResponse{
		ID:        userID,
		Email:     newUser.Email,
		UserName:  newUser.UserName,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		IsActive:  newUser.IsActive,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

// validateInput валидация входных данных
func (uc *RegisterUseCase) validateInput(input RegisterRequest) error {
	if err := uc.validateEmail(input.Email); err != nil {
		return err
	}

	if err := uc.validatePassword(input.Password); err != nil {
		return err
	}

	if input.Password != input.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	if strings.TrimSpace(input.UserName) == "" {
		return fmt.Errorf("username is required")
	}

	return nil
}

// validateEmail валидирует email
func (uc *RegisterUseCase) validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return user.ErrInvalidEmail
	}

	// Базовая проверка формата email
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return user.ErrInvalidEmail
	}

	return nil
}

func (uc *RegisterUseCase) validatePassword(password string) error {
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

// hashPassword хеширует пароль
func (uc *RegisterUseCase) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}
