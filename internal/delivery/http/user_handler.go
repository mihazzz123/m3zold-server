package http

import (
	"net/http"

	"github.com/mihazzz123/m3zold-server/internal/constants"
	"github.com/mihazzz123/m3zold-server/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	RegisterUC *user.RegisterUseCase
}

func NewUserHandler(registerUC *user.RegisterUseCase) *UserHandler {
	return &UserHandler{RegisterUC: registerUC}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input user.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	err := h.RegisterUC.Execute(c.Request.Context(), input)
	if err == constants.ErrEmailTaken {
		c.JSON(http.StatusConflict, gin.H{"error": "Email уже зарегистрирован"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Регистрация успешна"})
}
