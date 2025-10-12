package user

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

type LoginUseCase struct {
	Repo user.Repository
}

func NewLoginUseCase(repo user.Repository) *RegisterUseCase {
	return &RegisterUseCase{Repo: repo}
}

func (uc *LoginUseCase) generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (uc *LoginUseCase) generateVerificationToken() (string, error) {
	return uc.generateSecureToken(32)
}
