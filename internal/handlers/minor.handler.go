package handlers

import (
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/gin-gonic/gin"
)

func HandlePing() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := domain.
			PingResponseFormat{
			Success: true,
			Message: "pong",
		}
		c.JSON(http.StatusOK, &response)
	}
}
