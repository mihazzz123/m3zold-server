package services

import (
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

// UserFactory реализация фабрики пользователей
type UserFactory struct{}

// NewUserFactory создает новую UserFactory
func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

// CreateUser создает нового пользователя
func (f *UserFactory) CreateUser(id, email, userName, passwordHash, firstName, lastName string) *user.User {
	return user.NewUser(id, email, userName, passwordHash, firstName, lastName)
}
