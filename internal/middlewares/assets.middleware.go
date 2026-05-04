package middlewares

import (
	"log"
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func ValidatePublicKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		publicKey := c.Param("assetId")
		if publicKey == "" {
			log.Println("couldn't get assetId")
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Could not get the public key"})
			c.Abort()
			return
		}
		if _, ok := utils.ValidatePublicKey(publicKey); !ok {
			log.Println("Asset was invalid", publicKey)
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid public key"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CheckIfAssetIsPublic(cfg *configs.EnvData, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := utils.GetLogger(c)
		publicKey := c.Param("assetId")
		media, err := store.GetMedia(c.Request.Context(), logger, publicKey)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": http.StatusText(http.StatusNotFound)})
				c.Abort()
				return
			}
			log.Println("story", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			c.Abort()
			return
		}
		if !media.Public {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": http.StatusText(http.StatusForbidden)})
			c.Abort()
			return
		}
		c.Set("publicKey", publicKey)
		c.Set("media", media)
		c.Next()
	}
}
