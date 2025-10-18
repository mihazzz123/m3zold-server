// internal/domain/services/user_service.go
package services

import (
	"context"
	"fmt"

	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

type UserService struct {
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID возвращает пользователя по ID
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*user.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}
