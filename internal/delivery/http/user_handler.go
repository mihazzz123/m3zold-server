package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	"github.com/mihazzz123/m3zold-server/internal/infrastructure"
	userUS "github.com/mihazzz123/m3zold-server/internal/usecase/user"
)

type UserHandler struct {
	RegisterUC *userUS.RegisterUseCase
	AuthSrv    *infrastructure.AuthService
}

func NewUserHandler(registerUseCase *userUS.RegisterUseCase, authSrv *infrastructure.AuthService) *UserHandler {
	return &UserHandler{
		RegisterUC: registerUseCase,
		AuthSrv:    authSrv,
	}
}

// Register обработчик регистрации
func (h *UserHandler) Register(c *gin.Context) {
	var req userUS.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	response, err := h.RegisterUC.Execute(c.Request.Context(), req)
	if err != nil {
		status := http.StatusBadRequest

		switch err {
		case user.ErrEmailTaken:
			c.JSON(status, gin.H{"error": "Email already registered"})
		case user.ErrInvalidEmail, user.ErrWeakPassword:
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

func (h *UserHandler) Login(c *gin.Context) {
	// Пример: авторизация по userID
	userID := c.PostForm("user_id")
	userUIID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user_id"})
		return
	}
	token, err := h.AuthSrv.GenerateToken(cfg, userUIID)
	if err != nil {
		c.JSON(500, gin.H{"error": "token generation failed"})
		return
	}
	c.JSON(200, gin.H{"token": token})
}
