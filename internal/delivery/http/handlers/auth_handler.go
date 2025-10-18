package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mihazzz123/m3zold-server/internal/domain/services"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	"github.com/mihazzz123/m3zold-server/internal/usecase/auth"
)

type AuthHandler struct {
	AuthUseCase *auth.AuthUseCase
	UserService *services.UserService
}

func NewAuthHandler(authUseCase *auth.AuthUseCase, userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		AuthUseCase: authUseCase,
		UserService: userService,
	}
}

// Register обработчик регистрации
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	response, err := h.AuthUseCase.Register(c.Request.Context(), req)
	if err != nil {
		status := http.StatusBadRequest

		switch err {
		case user.ErrEmailTaken,
			user.ErrInvalidEmail,
			user.ErrPasswordConfirm,
			user.ErrPasswordRequired,
			user.ErrUserNameRequired,
			user.ErrEmailRequired,
			user.ErrWeakPassword,
			user.ErrUserNotFound,
			user.ErrInvalidCredentials:
			c.JSON(status, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registration successful",
		"user":    response,
	})
}

// Login обработчик входа
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	token, err := h.AuthUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"access_token":  token.Token,
			"token_type":    "Bearer",
			"expires_in":    900,                   // 15 minutes in seconds
			"refresh_token": "will_be_implemented", // TODO: вернуть реальный refresh token
		},
	})
}

// Logout обработчик выхода
func (h *AuthHandler) Logout(c *gin.Context) {
	// Используем ту же логику что в middleware для consistency
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authorization header required",
		})
		return
	}

	// Проверяем формат Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid authorization header format",
		})
		return
	}

	token := parts[1]

	if err := h.AuthUseCase.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// RefreshToken обработчик обновления токена
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	newToken, err := h.AuthUseCase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"access_token": newToken.Token,
			"token_type":   "Bearer",
			"expires_in":   900,
		},
	})
}

// ChangePassword обработчик смены пароля
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req auth.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Получаем userID из контекста (будет установлен middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	if err := h.AuthUseCase.ChangePassword(c.Request.Context(), userID.(string), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully",
	})
}

// GetCurrentUser возвращает данные текущего авторизованного пользователя
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Получаем userID из контекста (должен быть установлен в middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}

	// Получаем пользователя из сервиса
	user, err := h.UserService.GetUserByID(c.Request.Context(), userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get user data",
		})
		return
	}

	// Формируем ответ
	response := gin.H{
		"success": true,
		"data": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"userName":  user.UserName,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"isActive":  user.IsActive,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
}
