package user

import "errors"

var (
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrPasswordConfirm    = errors.New("password and confirmation do not match")
	ErrPasswordRequired   = errors.New("password is required")
	ErrUserNameRequired   = errors.New("username is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrWeakPassword       = errors.New("password does not meet security requirements")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
