package postgres

import (
	"context"

	"github.com/mihazzz123/m3zold-server/internal/domain/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{DB: db}
}

func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	return exists, err
}

func (r *userRepo) Create(ctx context.Context, u *user.User) error {
	_, err := r.DB.Exec(ctx,
		"INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3)",
		u.Email, u.PasswordHash, u.UserName)
	return err
}
