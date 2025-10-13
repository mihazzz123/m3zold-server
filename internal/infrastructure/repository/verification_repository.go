package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationEmailRepository struct {
	db *pgxpool.Pool
}

func NewVerificationEmailRepository(db *pgxpool.Pool) *VerificationEmailRepository {
	return &VerificationEmailRepository{db: db}
}

func (r *VerificationEmailRepository) CreateVerificationToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO m3zold_schema.verification_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create verification token: %w", err)
	}

	return nil
}

func (r *VerificationEmailRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	query := `
		SELECT user_id 
		FROM m3zold_schema.verification_tokens 
		WHERE token = $1 AND expires_at > $2 AND used = false
	`

	var userID string
	err := r.db.QueryRow(ctx, query, token, time.Now()).Scan(&userID)
	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("token not found or expired")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get user ID by token: %w", err)
	}

	return userID, nil
}

func (r *VerificationEmailRepository) DeleteToken(ctx context.Context, token string) error {
	query := `DELETE FROM m3zold_schema.verification_tokens WHERE token = $1`

	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

func (r *VerificationEmailRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	query := `UPDATE m3zold_schema.verification_tokens SET used = true WHERE token = $1`

	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}
