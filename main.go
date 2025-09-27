package main

import (
	"github.com/mihazzz123/m3zold-server/api/handlers"
	"github.com/mihazzz123/m3zold-server/api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger())

	r.GET("/auth/check", handlers.CheckAuth)

	r.Run(":8080") // слушаем порт
}
