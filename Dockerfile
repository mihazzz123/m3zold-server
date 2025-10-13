FROM golang:1.25.1-alpine

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем остальной код
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o m3zold-server ./cmd

EXPOSE 8080
CMD ["./m3zold-server"]