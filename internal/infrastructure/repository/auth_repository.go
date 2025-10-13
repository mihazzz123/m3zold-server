package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mihazzz123/m3zold-server/internal/domain/auth"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateToken(ctx context.Context, token *auth.Token) error {
	query := `
		INSERT INTO m3zold_schema.auth_tokens (
			id, user_id, token, token_type, expires_at, blacklisted, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		string(token.TokenType),
		token.ExpiresAt,
		token.Blacklisted,
		token.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	return nil
}

func (r *AuthRepository) GetToken(ctx context.Context, token string) (*auth.Token, error) {
	query := `
		SELECT id, user_id, token, token_type, expires_at, blacklisted, created_at
		FROM m3zold_schema.auth_tokens 
		WHERE token = $1
	`

	var dbToken auth.Token
	var tokenType string

	err := r.db.QueryRow(ctx, query, token).Scan(
		&dbToken.ID,
		&dbToken.UserID,
		&dbToken.Token,
		&tokenType,
		&dbToken.ExpiresAt,
		&dbToken.Blacklisted,
		&dbToken.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	dbToken.TokenType = auth.TokenType(tokenType)
	return &dbToken, nil
}

func (r *AuthRepository) BlacklistToken(ctx context.Context, token string) error {
	query := `UPDATE m3zold_schema.auth_tokens SET blacklisted = true WHERE token = $1`

	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

func (r *AuthRepository) DeleteUserTokens(ctx context.Context, userID string) error {
	query := `DELETE FROM m3zold_schema.auth_tokens WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user tokens: %w", err)
	}

	return nil
}

func (r *AuthRepository) CleanExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM m3zold_schema.auth_tokens WHERE expires_at < $1`

	_, err := r.db.Exec(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to clean expired tokens: %w", err)
	}

	return nil
}
