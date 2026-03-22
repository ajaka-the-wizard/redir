package middlewares

import (
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(mmap *memory.AuthMemoryMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie("sessionId")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": http.StatusText(http.StatusUnauthorized)})
			c.Abort()
			return
		}
		user, ok := mmap.GetUser(sessionId)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid sessionId"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
