package user

import "time"

type User struct {
	ID           string
	Email        string
	UserName     string
	PasswordHash string
	FirstName    string
	LastName     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// NewUser создает нового пользователя (фабричный метод)
func NewUser(id, email, userName, passwordHash, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		ID:           id,
		Email:        email,
		UserName:     userName,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    nil,
	}
}

// Activate активирует пользователя
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate деактивирует пользователя
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Delete помечает пользователя как удаленного
func (u *User) Delete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}

// UpdateProfile обновляет профиль пользователя
func (u *User) UpdateProfile(firstName, lastName string) {
	u.FirstName = firstName
	u.LastName = lastName
	u.UpdatedAt = time.Now()
}

// ChangeEmail изменяет email пользователя
func (u *User) ChangeEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}
