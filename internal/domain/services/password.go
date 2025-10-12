package services

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
}
