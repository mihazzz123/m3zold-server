package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mihazzz123/m3zold-server/internal/domain/user"
	userUC "github.com/mihazzz123/m3zold-server/internal/usecase/user"
)

type UserHandler struct {
	ProfileUseCase *userUC.ProfileUseCase
}

func NewUserHandler(profileUseCase *userUC.ProfileUseCase) *UserHandler {
	return &UserHandler{
		ProfileUseCase: profileUseCase,
	}
}

// GetProfile обработчик получения профиля пользователя
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	profile, err := h.ProfileUseCase.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// UpdateProfile обработчик обновления профиля
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	var req userUC.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	updatedProfile, err := h.ProfileUseCase.UpdateProfile(c.Request.Context(), userID.(string), req)
	if err != nil {
		status := http.StatusBadRequest
		if err == user.ErrEmailTaken {
			c.JSON(status, gin.H{
				"success": false,
				"error":   "Email already taken",
			})
			return
		}
		c.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedProfile,
		"message": "Profile updated successfully",
	})
}
