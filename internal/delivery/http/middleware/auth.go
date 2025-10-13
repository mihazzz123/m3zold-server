package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware добавляет user_id в контекст после аутентификации
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Здесь логика аутентификации (JWT, сессии и т.д.)
		// Для примера - берем из заголовка
		userID := c.GetHeader("X-User-ID")
		if userID != "" {
			ctx := context.WithValue(c.Request.Context(), "user_id", userID)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}
