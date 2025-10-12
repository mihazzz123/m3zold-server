package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	userUS "github.com/mihazzz123/m3zold-server/internal/usecase/user"
)

type UserHandler struct {
	RegisterUC *userUS.RegisterUseCase
}

func NewUserHandler(registerUseCase *userUS.RegisterUseCase) *UserHandler {
	return &UserHandler{
		RegisterUC: registerUseCase,
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
