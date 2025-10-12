package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

// UserRepository реализация репозитория пользователей для PostgreSQL
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepo(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	query := `
		INSERT INTO m3zold_schema.users (
			id, email, user_name, password_hash, 
			first_name, last_name, is_active, 
			created_at, updated_at, deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Конвертируем string ID в UUID для БД
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		userID,
		user.Email,
		user.UserName,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// ExistsByEmail проверяет существование пользователя по email
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM m3zold_schema.users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %v", err)
	}
	return exists, nil
}

// GetByEmail возвращает пользователя по email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, user_name, password_hash, first_name, last_name, 
			   is_active, created_at, updated_at, deleted_at
		FROM m3zold_schema.users 
		WHERE email = $1 AND deleted_at IS NULL
	`

	var (
		dbID   uuid.UUID
		dbUser user.User
	)

	err := r.db.QueryRow(ctx, query, email).Scan(
		&dbID,
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
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Конвертируем UUID обратно в string
	dbUser.ID = dbID.String()

	return &dbUser, nil
}

// GetByID возвращает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT id, email, user_name, password_hash, first_name, last_name, 
			   is_active, created_at, updated_at, deleted_at
		FROM m3zold_schema.users 
		WHERE id = $1 AND deleted_at IS NULL
	`

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var (
		dbID   uuid.UUID
		dbUser user.User
	)

	err = r.db.QueryRow(ctx, query, userID).Scan(
		&dbID,
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

	dbUser.ID = dbID.String()

	return &dbUser, nil
}

// Update обновляет данные пользователя
func (r *UserRepository) Update(ctx context.Context, user *user.User) error {
	query := `
		UPDATE m3zold_schema.users 
		SET email = $2, user_name = $3, first_name = $4, last_name = $5,
			is_active = $6, updated_at = $7, deleted_at = $8
		WHERE id = $1
	`

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		userID,
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
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE m3zold_schema.users 
		SET deleted_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL
	`

	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}

	return nil
}
