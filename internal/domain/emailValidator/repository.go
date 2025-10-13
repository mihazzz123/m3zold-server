package email

import (
	"context"
	"time"
)

// Repository интерфейс для репозитория email
type Repository interface {
	CreateVerificationToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	GetUserIDByToken(ctx context.Context, token string) (string, error)
	DeleteToken(ctx context.Context, token string) error
	MarkTokenAsUsed(ctx context.Context, token string) error
}
