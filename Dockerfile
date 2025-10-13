FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Копируем зависимости first для лучшего кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка бинарника с оптимизациями
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o m3zold-server ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /app/m3zold-server .
COPY --from=builder /app/migrations ./migrations

# Создаем non-root пользователя для безопасности
RUN adduser -D -s /bin/sh appuser
RUN chown -R appuser:appuser /root/
USER appuser

EXPOSE 8080

CMD ["./m3zold-server"]