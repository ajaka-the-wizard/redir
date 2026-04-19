package middlewares

import (
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckIfAssetIsPublic(cfg *configs.EnvData, pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		assetId := c.Param("assetId")
		if assetId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Could not get the asset id"})
			c.Abort()
			return
		}
		_, ok := utils.ValidateAssetId(assetId)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid asset id"})
			c.Abort()
			return
		}
		normalised := cfg.DATA_GET_PATH + "/" + assetId
		media, err := repository.GetMedia(c.Request.Context(), pool, normalised)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": http.StatusText(http.StatusNotFound)})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			c.Abort()
			return
		}
		if !media.Public {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": http.StatusText(http.StatusForbidden)})
			c.Abort()
			return
		}
		c.Set("publicKey", normalised)
		c.Set("media", &media)
		c.Next()
	}
}
