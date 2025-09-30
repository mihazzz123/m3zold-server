package user

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/constants"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email    string
	Password string
	Name     string
}

type RegisterUseCase struct {
	Repo user.Repository
}

func NewRegisterUseCase(repo user.Repository) *RegisterUseCase {
	return &RegisterUseCase{Repo: repo}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) error {
	exists, err := uc.Repo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if exists {
		return constants.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &user.User{
		Email:        input.Email,
		PasswordHash: string(hash),
		Name:         input.Name,
	}

	return uc.Repo.Create(ctx, user)
}
