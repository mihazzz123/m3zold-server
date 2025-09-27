FROM golang:1.25.1-alpine

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o go-laeg ./src

EXPOSE 8080
CMD ["./go-laeg"]
