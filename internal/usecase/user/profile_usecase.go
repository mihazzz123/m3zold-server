package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mihazzz123/m3zold-server/internal/domain/services"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

type ProfileUseCase struct {
	userRepo       user.Repository
	emailValidator services.EmailValidatorService
}

func NewProfileUseCase(userRepo user.Repository, emailValidator services.EmailValidatorService) *ProfileUseCase {
	return &ProfileUseCase{
		userRepo:       userRepo,
		emailValidator: emailValidator,
	}
}

// GetProfileResponse DTO для ответа профиля
type GetProfileResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateProfileRequest DTO для обновления профиля
type UpdateProfileRequest struct {
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateProfileResponse DTO для ответа обновления профиля
type UpdateProfileResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetProfile получает профиль пользователя
func (uc *ProfileUseCase) GetProfile(ctx context.Context, userID string) (*GetProfileResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &GetProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// UpdateProfile обновляет профиль пользователя
func (uc *ProfileUseCase) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*UpdateProfileResponse, error) {
	// Получаем текущего пользователя
	currentUser, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Валидация email если он изменился
	if req.Email != "" && req.Email != currentUser.Email {
		normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

		if err := uc.emailValidator.Validate(normalizedEmail); err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}

		// Проверяем что email не занят другим пользователем
		exists, err := uc.userRepo.ExistsByEmail(ctx, normalizedEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if exists {
			return nil, user.ErrEmailTaken
		}

		currentUser.Email = normalizedEmail
	}

	// Обновляем остальные поля если они предоставлены
	if req.UserName != "" {
		currentUser.UserName = strings.TrimSpace(req.UserName)
	}

	if req.FirstName != "" {
		currentUser.FirstName = strings.TrimSpace(req.FirstName)
	}

	if req.LastName != "" {
		currentUser.LastName = strings.TrimSpace(req.LastName)
	}

	// Обновляем timestamp
	currentUser.UpdatedAt = time.Now()

	// Сохраняем изменения
	if err := uc.userRepo.Update(ctx, currentUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &UpdateProfileResponse{
		ID:        currentUser.ID,
		Email:     currentUser.Email,
		UserName:  currentUser.UserName,
		FirstName: currentUser.FirstName,
		LastName:  currentUser.LastName,
		UpdatedAt: currentUser.UpdatedAt,
	}, nil
}

// ChangeEmail меняет email пользователя
func (uc *ProfileUseCase) ChangeEmail(ctx context.Context, userID, newEmail string) error {
	// Получаем пользователя
	curUser, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Нормализуем и валидируем новый email
	normalizedEmail := strings.ToLower(strings.TrimSpace(newEmail))
	if err := uc.emailValidator.Validate(normalizedEmail); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	// Проверяем что email не занят
	exists, err := uc.userRepo.ExistsByEmail(ctx, normalizedEmail)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return user.ErrEmailTaken
	}

	// Обновляем email
	curUser.Email = normalizedEmail
	curUser.UpdatedAt = time.Now()

	return uc.userRepo.Update(ctx, curUser)
}

// DeactivateProfile деактивирует профиль пользователя
func (uc *ProfileUseCase) DeactivateProfile(ctx context.Context, userID string) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.Deactivate()
	user.UpdatedAt = time.Now()

	return uc.userRepo.Update(ctx, user)
}
