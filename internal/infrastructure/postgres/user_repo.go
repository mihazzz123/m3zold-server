package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{DB: db}
}

// Create создает нового пользователя
func (r *userRepo) Create(ctx context.Context, user *user.User) error {
	query := `
        INSERT INTO m3zold_schema.users (
            id,
            email,
            user_name,
            password_hash,
            first_name,
            last_name,
            is_active,
            created_at,
            updated_at,
            deleted_at,
			is_verified
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

	_, err := r.DB.Exec(ctx, query,
		user.ID,
		user.Email,
		user.UserName,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
		user.IsVerified,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// ExistsByEmail проверяет существование пользователя по email
func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM m3zold_schema.users WHERE email = $1)`
	var exists bool
	err := r.DB.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// GetByEmail возвращает пользователя по email
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, user_name, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at, is_verified
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
		logrus.Errorf("failed to get user by email: %s", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &u, nil
}

// GetByID возвращает пользователя по ID
func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, email, user_name, password_hash, first_name, last_name, 
			   is_active, created_at, updated_at, deleted_at
		FROM m3zold_schema.users 
		WHERE id = $1 AND deleted_at IS NULL
	`

	var dbUser user.User

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&dbUser.ID,
		&dbUser.Email,
		&dbUser.UserName,
		&dbUser.PasswordHash,
		&dbUser.FirstName,
		&dbUser.LastName,
		&dbUser.IsActive,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
		&dbUser.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &dbUser, nil
}

// Update обновляет данные пользователя
func (r *userRepo) Update(ctx context.Context, user *user.User) error {
	query := `
		UPDATE m3zold_schema.users 
		SET email = $2, user_name = $3, first_name = $4, last_name = $5,
			is_active = $6, updated_at = $7, deleted_at = $8
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query,
		user.ID,
		user.Email,
		user.UserName,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.UpdatedAt,
		user.DeletedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete помечает пользователя как удаленного
func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE m3zold_schema.users 
		SET deleted_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}

	return nil
}
