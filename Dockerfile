FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Копируем зависимости first для лучшего кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o m3zold-server ./cmd

FROM alpine:latest

# Обновляем apk репозитории и устанавливаем пакеты с ретраями
RUN apk update --no-cache && \
    apk add --no-cache ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /app/m3zold-server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./m3zold-server"]