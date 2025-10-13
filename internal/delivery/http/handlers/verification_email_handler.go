package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	email "github.com/mihazzz123/m3zold-server/internal/usecase/verification_email"
)

type VerificationEmailHandler struct {
	emailUseCase *email.VerificationEmailUseCase
}

func NewVerificationEmailHandler(emailUseCase *email.VerificationEmailUseCase) *VerificationEmailHandler {
	return &VerificationEmailHandler{
		emailUseCase: emailUseCase,
	}
}

// VerifyEmail обработчик верификации email
// @Summary Verify email
// @Description Verify user's email address with token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]interface{} "Email verified successfully"
// @Failure 400 {object} map[string]interface{} "Invalid token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/verify-email [post]
func (h *VerificationEmailHandler) VerifyEmail(c *gin.Context) {
	var req email.VerifyEmailRequest

	// Получаем token из query параметра
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Verification token is required",
		})
		return
	}

	req.Token = token

	// Верифицируем email
	userID, err := h.emailUseCase.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid or expired verification token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email verified successfully",
		"user_id": userID,
	})
}

// ResendVerification обработчик повторной отправки верификации
// @Summary Resend verification email
// @Description Resend email verification link
// @Tags auth
// @Accept json
// @Produce json
// @Param request body email.ResendVerificationRequest true "Resend verification request"
// @Success 200 {object} map[string]interface{} "Verification email sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/resend-verification [post]
func (h *VerificationEmailHandler) ResendVerification(c *gin.Context) {
	var req email.ResendVerificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if err := h.emailUseCase.ResendVerification(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to resend verification email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification email sent successfully",
	})
}

// SendWelcomeEmail обработчик отправки приветственного email
// @Summary Send welcome email
// @Description Send welcome email to user
// @Tags users
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Param name query string true "User name"
// @Success 200 {object} map[string]interface{} "Welcome email sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /email/welcome [post]
func (h *VerificationEmailHandler) SendWelcomeEmail(c *gin.Context) {
	email := c.Query("email")
	name := c.Query("name")

	if email == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Email and name are required",
		})
		return
	}

	if err := h.emailUseCase.SendWelcomeEmail(c.Request.Context(), email, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to send welcome email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Welcome email sent successfully",
	})
}
