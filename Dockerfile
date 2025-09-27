FROM golang:1.25.1-alpine

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o m3zold-server ./src

EXPOSE 8080
CMD ["./m3zold-server"]
