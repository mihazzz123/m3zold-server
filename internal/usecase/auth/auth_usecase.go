package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mihazzz123/m3zold-server/internal/domain/auth"
	"github.com/mihazzz123/m3zold-server/internal/domain/services"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
)

type AuthUseCase struct {
	authRepo        auth.Repository
	userRepo        user.Repository
	passwordService services.PasswordService
	tokenService    services.TokenService
	jwtService      services.JWTService
	idService       services.IDService
	emailValidator  services.EmailValidatorService
	userFactory     services.UserFactory
}

func NewAuthUseCase(
	authRepo auth.Repository,
	userRepo user.Repository,
	passwordService services.PasswordService,
	tokenService services.TokenService,
	jwtService services.JWTService,
	emailValidator services.EmailValidatorService,
	idService services.IDService,
	userFactory services.UserFactory,
) *AuthUseCase {
	return &AuthUseCase{
		authRepo:        authRepo,
		userRepo:        userRepo,
		passwordService: passwordService,
		tokenService:    tokenService,
		jwtService:      jwtService,
		idService:       idService,
		emailValidator:  emailValidator,
		userFactory:     userFactory,
	}
}

type RegisterRequest struct {
	Email           string `json:"email"`
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
}

type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
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

func (uc *AuthUseCase) Register(ctx context.Context, input RegisterRequest) (*RegisterResponse, error) {
	// Валидация входных данных
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Нормализация email
	email := strings.ToLower(strings.TrimSpace(input.Email))

	// Валидация email через domain service
	if err := uc.emailValidator.Validate(email); err != nil {
		return nil, err
	}

	// Проверка существования пользователя
	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if exists {
		return nil, user.ErrEmailTaken
	}

	// Хеширование пароля через domain service
	passwordHash, err := uc.passwordService.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("password processing error: %w", err)
	}

	// Генерация ID через domain service
	userID := uc.idService.Generate()

	// Создание пользователя через фабрику
	newUser := uc.userFactory.CreateUser(
		userID,
		email,
		input.UserName,
		passwordHash,
		input.FirstName,
		input.LastName,
	)

	// Сохранение пользователя
	if err = uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Возвращаем ответ
	return &RegisterResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		UserName:  newUser.UserName,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		IsActive:  newUser.IsActive,
	}, nil
}

func (uc *AuthUseCase) validateInput(input RegisterRequest) error {
	if strings.TrimSpace(input.Email) == "" {
		return user.ErrEmailRequired
	}

	if strings.TrimSpace(input.UserName) == "" {
		return user.ErrUserNameRequired
	}

	if strings.TrimSpace(input.Password) == "" {
		return user.ErrPasswordRequired
	}

	if input.Password != input.ConfirmPassword {
		return user.ErrPasswordConfirm
	}

	if len(input.Password) < 8 {
		return user.ErrWeakPassword
	}

	return nil
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
