package auth

import (
	"context"
	"time"
)

// TokenType тип токена
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// Token структура токена
type Token struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Token       string    `json:"token"`
	TokenType   TokenType `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	Blacklisted bool      `json:"blacklisted"`
	CreatedAt   time.Time `json:"created_at"`
}

// Claims JWT claims
type Claims struct {
	UserID   string    `json:"user_id"`
	Email    string    `json:"email"`
	UserName string    `json:"user_name"`
	Exp      time.Time `json:"exp"`
}

// UseCase интерфейс для use case аутентификации
type UseCase interface {
	Login(ctx context.Context, email, password string) (*Token, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*Token, error)
	ValidateToken(ctx context.Context, token string) (*Claims, error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
}
