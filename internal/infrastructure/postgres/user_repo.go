package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{DB: db}
}

func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM m3zold_schema.users WHERE email = $1)`
	var exists bool
	err := r.DB.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

func (r *userRepo) Create(ctx context.Context, user *user.User) error {
	query := `
		INSERT INTO m3zold_schema.users (
			id,
			email,
			password_hash,
			user_name,
			first_name,
			second_name,
			created_at,
			updated_at,
			deleted_at,
			is_active,
			is_verified		
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.DB.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.UserName,
		user.FirsName,
		user.SecondName,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedDt,
		user.IsActive,
		user.IsVerified,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at, is_active, is_verified
		FROM m3zold_schema.users WHERE email = $1
	`

	var u user.User
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsActive,
		&u.IsVerified,
	)

	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		logrus.Errorf("failed to get user by email: %w", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &u, nil
}

func (r *userRepo) Update(ctx context.Context, user *user.User) error {
	query := `
		UPDATE m3zold_schema.users
		SET email = $1,	 password_hash = $2, user_name = $3,
		    first_name = $4, second_name = $5,
			updated_at = $6, deleted_at = $7,
			is_active = $8, is_verified = $9
		WHERE id = $10
	`
	_, err := r.DB.Exec(ctx, query,
		user.Email,
		user.PasswordHash, user.UserName,
		user.FirsName, user.SecondName,
		user.UpdatedAt, user.DeletedDt,
		user.IsActive, user.IsVerified,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
