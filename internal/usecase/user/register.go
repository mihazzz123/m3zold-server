package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mihazzz123/m3zold-server/internal/constants"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

type RegisterUseCase struct {
	Repo user.Repository
	// EmailSender EmailSender
}

// RegisterRequest содержит данные для регистрации пользователя
type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	UserName        string `json:"user_name"`
}

func NewRegisterUseCase(repo user.Repository) *RegisterUseCase {
	return &RegisterUseCase{Repo: repo}
}

// Execute выполняет регистрацию пользователя
func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterRequest) (*user.User, error) {
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

	return nil
}
