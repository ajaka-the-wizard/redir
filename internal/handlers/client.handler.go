package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/gin-gonic/gin"
)

func HandleClientPing() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get("client")
		client, ok := val.(*models.Product)
		if !ok {
			log.Println("bbb")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Client %v recognised, Proxy is ready for you.", client.ProductId)})
	}
}
