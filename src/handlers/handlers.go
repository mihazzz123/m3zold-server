package handlers

import "github.com/gin-gonic/gin"

func CheckAuth(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok", "message": "Auth check passed!"})
}
