package auth

import "context"

// Repository интерфейс для репозитория аутентификации
type Repository interface {
	CreateToken(ctx context.Context, token *Token) error
	GetToken(ctx context.Context, token string) (*Token, error)
	BlacklistToken(ctx context.Context, token string) error
	DeleteUserTokens(ctx context.Context, userID string) error
	CleanExpiredTokens(ctx context.Context) error
}
