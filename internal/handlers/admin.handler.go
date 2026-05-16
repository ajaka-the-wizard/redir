package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlePing() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "pong"})
	}
}
