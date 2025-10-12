package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationRepository struct {
	db *pgxpool.Pool
}

func NewVerificationRepository(db *pgxpool.Pool) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) CreateVerificationToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO verification_tokens (token, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query, token, userID, expiresAt, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create verification token: %w", err)
	}

	return nil
}

func (r *VerificationRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	query := `
		SELECT user_id 
		FROM verification_tokens 
		WHERE token = $1 AND expires_at > $2 AND used = false
	`

	var userID string
	err := r.db.QueryRow(ctx, query, token, time.Now()).Scan(&userID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("token not found or expired")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get user ID by token: %w", err)
	}

	return userID, nil
}

func (r *VerificationRepository) DeleteToken(ctx context.Context, token string) error {
	query := `DELETE FROM verification_tokens WHERE token = $1`

	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

func (r *VerificationRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	query := `UPDATE verification_tokens SET used = true WHERE token = $1`

	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}
