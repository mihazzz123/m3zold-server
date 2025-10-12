package services

import (
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

// PasswordHasher интерфейс для хеширования паролей
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) bool
}

// IDGenerator интерфейс для генерации ID
type IDGenerator interface {
	Generate() string
}

// EmailValidator интерфейс для валидации email
type EmailValidator interface {
	Validate(email string) error
}

// TokenGenerator интерфейс для генерации токенов
type TokenGenerator interface {
	GenerateToken() (string, error)
	GenerateTokenWithLength(length int) (string, error)
}

// UserFactory интерфейс для создания пользователей
type UserFactory interface {
	CreateUser(id, email, userName, passwordHash, firstName, lastName string) *user.User
}
