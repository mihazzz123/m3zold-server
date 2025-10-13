FROM golang:1.25.1-alpine

WORKDIR /app

# Копируем vendor папку и модули
COPY vendor ./vendor
COPY go.mod go.sum ./

# Копируем .env файл в контейнер
COPY .env ./

# Копируем остальной код
COPY . .

# Сборка бинарника (используем vendor папку)
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o m3zold-server ./cmd

EXPOSE 8080
CMD ["./m3zold-server"]