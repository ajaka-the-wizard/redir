package handlers

import (
	"fmt"
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
)

func HandleClientPing() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := utils.GetLogger(c)
		logger.Info("ping attempt from client")
		val, _ := c.Get("client")
		client, ok := val.(*models.Product)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Client %v recognised, Proxy is ready for you.", client.ProductId)})
		logger.Info("Client ping ponged", "productId", client.ProductId)
	}
}
