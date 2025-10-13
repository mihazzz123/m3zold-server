package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ContextMiddleware добавляет request_id и логгер в контекст
func ContextMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Генерируем request ID
		requestID := uuid.New().String()

		// Создаем контекст с таймаутом и request_id
		ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
		defer cancel()

		// Добавляем значения в контекст
		ctx = context.WithValue(ctx, "request_id", requestID)
		ctx = context.WithValue(ctx, "ip_address", c.ClientIP())
		ctx = context.WithValue(ctx, "user_agent", c.Request.UserAgent())

		// Создаем логгер для этого запроса
		requestLogger := logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		ctx = context.WithValue(ctx, "logger", requestLogger)

		// Заменяем контекст в запросе
		c.Request = c.Request.WithContext(ctx)

		// Добавляем request_id в заголовки ответа
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
