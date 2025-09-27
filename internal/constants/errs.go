package constants

import "errors"

var (
	// ErrEmailTaken ошибка, возвращаемая при попытке регистрации с уже существующим email
	ErrEmailTaken error = errors.New("email already taken")
)
