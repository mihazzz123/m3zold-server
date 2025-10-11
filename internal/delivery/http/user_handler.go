package http

import (
	"encoding/json"
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
	// Ограничение частоты запросов (реализуйте отдельно)
	// if err := h.rateLimiter.Check(r); err != nil {
	// http.Error(w, "Too many requests", http.StatusTooManyRequests)
	// return
	// }

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if errors := req.Validate(); len(errors) > 0 {
		respondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"errors":  errors,
		})
		return
	}

	/**

		    user, err := h.authService.Register(r.Context(), &req)
	    if err != nil {
	        // Обобщённое сообщение об ошибке для безопасности
	        status := http.StatusBadRequest
	        if strings.Contains(err.Error(), "database error") {
	            status = http.StatusInternalServerError
	        }

	        respondWithJSON(w, status, map[string]interface{}{
	            "success": false,
	            "message": "Registration failed",
	        })
	        return
	    }

	    respondWithJSON(w, http.StatusCreated, map[string]interface{}{
	        "success": true,
	        "message": "Registration successful. Please check your email for verification.",
	        "user":    user,
	    })
	*/

	newUser, err := h.RegisterUC.Execute(c.Request.Context(), input)
	if err == constants.ErrEmailTaken {
		c.JSON(http.StatusConflict, gin.H{"error": "Email уже зарегистрирован"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}
	if newUser == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания нового пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Регистрация успешна"})
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
