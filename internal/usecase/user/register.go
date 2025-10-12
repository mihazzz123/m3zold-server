package userusecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/mihazzz123/m3zold-server/internal/domain/services"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

type RegisterUseCase struct {
	userRepo       user.Repository
	passwordHasher services.PasswordHasher
	idGenerator    services.IDGenerator
	emailValidator services.EmailValidator
}

func NewRegisterUseCase(
	userRepo user.Repository,
	passwordHasher services.PasswordHasher,
	idGenerator services.IDGenerator,
	emailValidator services.EmailValidator,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		idGenerator:    idGenerator,
		emailValidator: emailValidator,
	}
}

type RegisterRequest struct {
	Email           string `json:"email"`
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
}

type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterRequest) (*RegisterResponse, error) {
	// Валидация входных данных
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Нормализация email
	email := strings.ToLower(strings.TrimSpace(input.Email))

	// Валидация email через domain service
	if err := uc.emailValidator.Validate(email); err != nil {
		return nil, err
	}

	// Проверка существования пользователя
	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if exists {
		return nil, user.ErrEmailTaken
	}

	// Хеширование пароля через domain service
	passwordHash, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("password processing error: %w", err)
	}

	// Генерация ID через domain service
	userID := uc.idGenerator.Generate()

	// Создание пользователя через фабричный метод
	newUser := user.NewUser(
		userID,
		email,
		input.UserName,
		passwordHash,
		input.FirstName,
		input.LastName,
	)

	// Сохранение пользователя
	if err = uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Возвращаем ответ
	return &RegisterResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		UserName:  newUser.UserName,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		IsActive:  newUser.IsActive,
	}, nil
}

func (uc *RegisterUseCase) validateInput(input RegisterRequest) error {
	if strings.TrimSpace(input.Email) == "" {
		return fmt.Errorf("email is required")
	}

	if strings.TrimSpace(input.UserName) == "" {
		return fmt.Errorf("username is required")
	}

	if strings.TrimSpace(input.Password) == "" {
		return fmt.Errorf("password is required")
	}

	if input.Password != input.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	if len(input.Password) < 8 {
		return user.ErrWeakPassword
	}

	return nil
}
