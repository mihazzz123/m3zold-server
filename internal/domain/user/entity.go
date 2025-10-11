package user

import "time"

type User struct {
	ID           string     `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	UserName     string     `json:"user_name" db:"user_name"`
	FirsName     string     `json:"first_name" db:"first_name"`
	SecondName   string     `json:"second_name" db:"second_name"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedDt    *time.Time `json:"deleted_at" db:"deleted_at"`
	IsActive     bool       `json:"is_active" db:"is_active"`
}
