package email

import (
	"context"
	"fmt"
	"time"

	"github.com/mihazzz123/m3zold-server/internal/domain/emailVerification"
	"github.com/mihazzz123/m3zold-server/internal/domain/services"
)

type VerificationEmailUseCase struct {
	emailRepo      email.Repository
	tokenGenerator services.TokenService
}

func NewVerificationEmailUseCase(emailRepo email.Repository, tokenGenerator services.TokenService) *VerificationEmailUseCase {
	return &VerificationEmailUseCase{
		emailRepo:      emailRepo,
		tokenGenerator: tokenGenerator,
	}
}

// SendVerificationRequest DTO для запроса отправки верификации
type Request struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmailRequest DTO для запроса верификации email
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmailResponse DTO для ответа верификации
type VerifyEmailResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// ResendVerificationRequest DTO для повторной отправки верификации
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (uc *VerificationEmailUseCase) SendVerificationEmail(ctx context.Context, userID, userEmail string) error {
	// Генерация токена верификации
	token, err := uc.tokenGenerator.GenerateToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Сохранение токена в БД (действителен 24 часа)
	expiresAt := time.Now().Add(24 * time.Hour)
	if err := uc.emailRepo.CreateVerificationToken(ctx, userID, token, expiresAt); err != nil {
		return fmt.Errorf("failed to save verification token: %w", err)
	}

	// TODO: Реальная отправка email
	// Здесь будет интеграция с email сервисом (SMTP, SendGrid, etc.)
	fmt.Printf("📧 Verification email would be sent to: %s, token: %s\n", userEmail, token)

	return nil
}

func (uc *VerificationEmailUseCase) VerifyEmail(ctx context.Context, token string) (string, error) {
	userID, err := uc.emailRepo.GetUserIDByToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("invalid or expired token: %w", err)
	}

	// Помечаем токен как использованный
	if err := uc.emailRepo.MarkTokenAsUsed(ctx, token); err != nil {
		// Логируем ошибку, но не прерываем процесс верификации
		fmt.Printf("Failed to mark token as used: %v\n", err)
	}

	return userID, nil
}

func (uc *VerificationEmailUseCase) SendWelcomeEmail(ctx context.Context, email, name string) error {
	// TODO: Реальная отправка приветственного email
	fmt.Printf("📧 Welcome email would be sent to: %s, name: %s\n", email, name)
	return nil
}

func (uc *VerificationEmailUseCase) SendPasswordResetEmail(ctx context.Context, email, token string) error {
	// TODO: Реальная отправка email сброса пароля
	fmt.Printf("📧 Password reset email would be sent to: %s, token: %s\n", email, token)
	return nil
}

// ResendVerification повторно отправляет email верификации
func (uc *VerificationEmailUseCase) ResendVerification(ctx context.Context, email string) error {
	// TODO: Реализовать логику поиска пользователя по email и повторной отправки
	// Пока заглушка
	fmt.Printf("📧 Resending verification email to: %s\n", email)
	return nil
}
