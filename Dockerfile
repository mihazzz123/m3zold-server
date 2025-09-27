FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go build -o go-laeg main.go

EXPOSE 8080
CMD ["./go-laeg"]
