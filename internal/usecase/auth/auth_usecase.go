package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/mihazzz123/m3zold-server/internal/domain/auth"
	"github.com/mihazzz123/m3zold-server/internal/domain/services"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

type AuthUseCase struct {
	userRepo        user.Repository
	authRepo        auth.Repository
	passwordService services.PasswordService
	tokenService    services.TokenService
	jwtService      services.JWTService
	idService       services.IDService
}

func NewAuthUseCase(
	userRepo user.Repository,
	authRepo auth.Repository,
	passwordService services.PasswordService,
	tokenService services.TokenService,
	jwtService services.JWTService,
	idService services.IDService,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:        userRepo,
		authRepo:        authRepo,
		passwordService: passwordService,
		tokenService:    tokenService,
		jwtService:      jwtService,
		idService:       idService,
	}
}

// LoginRequest DTO для входа
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse DTO для ответа входа
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		UserName string `json:"user_name"`
	} `json:"user"`
}

// RefreshRequest DTO для обновления токена
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest DTO для смены пароля
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*auth.Token, error) {
	// Находим пользователя по email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Проверяем пароль
	if !uc.passwordService.Verify(user.PasswordHash, password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Проверяем что пользователь активен
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Генерируем access token
	accessToken, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.UserName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Генерируем refresh token
	refreshToken, err := uc.tokenService.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Сохраняем refresh token в БД
	refreshTokenEntity := &auth.Token{
		ID:        uc.idService.Generate(),
		UserID:    user.ID,
		Token:     refreshToken,
		TokenType: auth.TokenTypeRefresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 дней
		CreatedAt: time.Now(),
	}

	if err := uc.authRepo.CreateToken(ctx, refreshTokenEntity); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &auth.Token{
		ID:        uc.idService.Generate(),
		UserID:    user.ID,
		Token:     accessToken,
		TokenType: auth.TokenTypeAccess,
		ExpiresAt: time.Now().Add(15 * time.Minute), // 15 минут
		CreatedAt: time.Now(),
	}, nil
}

func (uc *AuthUseCase) Logout(ctx context.Context, token string) error {
	return uc.authRepo.BlacklistToken(ctx, token)
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*auth.Token, error) {
	// Находим refresh token в БД
	token, err := uc.authRepo.GetToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Проверяем не истек ли токен
	if time.Now().After(token.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Проверяем не заблокирован ли токен
	if token.Blacklisted {
		return nil, fmt.Errorf("token blacklisted")
	}

	// Находим пользователя
	user, err := uc.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Генерируем новый access token
	accessToken, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.UserName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &auth.Token{
		ID:        uc.idService.Generate(),
		UserID:    user.ID,
		Token:     accessToken,
		TokenType: auth.TokenTypeAccess,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
	}, nil
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
	return uc.jwtService.ValidateToken(token)
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Находим пользователя
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Проверяем старый пароль
	if !uc.passwordService.Verify(user.PasswordHash, oldPassword) {
		return fmt.Errorf("invalid old password")
	}

	// Хешируем новый пароль
	newPasswordHash, err := uc.passwordService.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Обновляем пароль
	user.PasswordHash = newPasswordHash
	user.UpdatedAt = time.Now()

	// TODO: Добавить метод Update в user repository
	// if err := uc.userRepo.Update(ctx, user); err != nil {
	//     return fmt.Errorf("failed to update password: %w", err)
	// }

	return nil
}
