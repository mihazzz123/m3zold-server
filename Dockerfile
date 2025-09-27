FROM golang:1.25.1-alpine

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod 
RUN go mod vendor

# Копируем остальной код
COPY . .

# Сборка бинарника
RUN go build -v -o m3zold-server ./src

EXPOSE 8080
CMD ["./m3zold-server"]
