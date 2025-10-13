package services

import (
	"github.com/mihazzz123/m3zold-server/internal/domain/auth"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

// PasswordService интерфейс для хеширования паролей
type PasswordService interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) bool
}

// IDService интерфейс для генерации ID
type IDService interface {
	Generate() string
}

// EmailValidatorService интерфейс для валидации email
type EmailValidatorService interface {
	Validate(email string) error
	Normalize(email string) string
}

// TokenService интерфейс для генерации токенов
type TokenService interface {
	GenerateToken() (string, error)
	GenerateTokenWithLength(length int) (string, error)
}

// UserFactory интерфейс для создания пользователей
type UserFactory interface {
	CreateUser(id, email, userName, passwordHash, firstName, lastName string) *user.User
}

// JWTService интерфейс для работы с JWT токенами
type JWTService interface {
	GenerateToken(userID, email, userName string) (string, error)
	ValidateToken(token string) (*auth.Claims, error)
}
