package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure"
)

type RegisterUseCase struct {
	Repo        user.Repository
	PasswordSvc infrastructure.PasswordService
}

func NewRegisterUseCase(repo user.Repository, passwordSvc infrastructure.PasswordService) *RegisterUseCase {
	return &RegisterUseCase{Repo: repo, PasswordSvc: passwordSvc}
}

// RegisterRequest DTO для регистрации
type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// RegisterResponse DTO ответа
type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
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
	passwordHash, err := uc.PasswordSvc.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("password processing error: %w", err)
	}

	// Создание пользователя
	newUser := &user.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true, // или false если требуется верификация
		IsVerified:   false,
	}

	// Сохранение в репозиторий
	if err = uc.Repo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Возвращаем ответ без чувствительных данных
	return &RegisterResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		IsActive:  newUser.IsActive,
	}, nil
}

// validateInput валидация входных данных
func (uc *RegisterUseCase) validateInput(input RegisterRequest) error {
	if err := uc.PasswordSvc.ValidateEmail(input.Email); err != nil {
		return err
	}

	if err := uc.PasswordSvc.ValidatePassword(input.Password); err != nil {
		return err
	}

	if input.Password != input.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	return nil
}
