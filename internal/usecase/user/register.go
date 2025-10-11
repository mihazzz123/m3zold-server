package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	"github.com/sirupsen/logrus"
)

type RegisterUseCase struct {
	Repo user.Repository
	// EmailSender EmailSender
}

func NewRegisterUseCase(repo user.Repository) *RegisterUseCase {
	return &RegisterUseCase{Repo: repo}
}

// Execute выполняет регистрацию пользователя
func (uc *RegisterUseCase) Execute(ctx context.Context, req *user.RegisterRequest) (*user.User, error) {
	// Валидация входных данных
	if errors := Validate(req); len(errors) > 0 {
		logrus.Errorf("validation errors: %v", errors)
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	exists, err := uc.Repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, user.ErrEmailTaken
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
