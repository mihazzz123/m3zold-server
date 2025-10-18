// internal/delivery/http/middleware/auth_middleware.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header required",
			})
			c.Abort()
			return
		}

		// Проверяем формат "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Парсим JWT токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// TODO: Заменить на ваш секретный ключ
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token",
			})
			c.Abort()
			return
		}

		// Извлекаем claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Устанавливаем userID в контекст
			if userID, exists := claims["user_id"]; exists {
				if userIDStr, ok := userID.(string); ok {
					c.Set("userID", userIDStr)
				}
			}

			// Можно также установить другие данные из токена
			if email, exists := claims["email"]; exists {
				if emailStr, ok := email.(string); ok {
					c.Set("email", emailStr)
				}
			}

			if userName, exists := claims["user_name"]; exists {
				if userNameStr, ok := userName.(string); ok {
					c.Set("userName", userNameStr)
				}
			}
		}

		c.Next()
	}
}
