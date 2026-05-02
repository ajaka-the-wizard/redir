package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func AuthMiddleware(store *store.Store, cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie("sessionId")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": http.StatusText(http.StatusUnauthorized)})
			c.Abort()
			return
		}
		user, ok := store.GetUser(c.Request.Context(), sessionId)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired sessionId"})
			c.Abort()
			return
		}
		c.Set("sessionId", sessionId)
		c.Set("user", user)
		c.Next()
	}
}

func CheckAndValidateClientKeys(cfg *configs.EnvData, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		pId := c.GetHeader("X-Product")
		pKey := c.GetHeader("Authorization")
		k, ok := strings.CutPrefix(pKey, "Bearer ")
		if !ok || pId == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Auth token or clientId"})
			c.Abort()
			return
		}
		pIdI, err := strconv.Atoi(pId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
			c.Abort()
			return
		}
		product, err := store.GetProductById(c.Request.Context(), pIdI)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Product with id of %d not found", pIdI)})
				c.Abort()
				return
			}
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			c.Abort()
			return
		}
		err = utils.VerifyMultipStepHash(k, product.PrivateKey)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid key given"})
			c.Abort()
			return
		}
		c.Set("product", product)
		c.Next()
	}
}

func CanThisUserAlterThisProduct(cfg *configs.EnvData, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := utils.GetUser(c)
		if !ok || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			c.Abort()
			return
		}
		pId := c.Param("id")
		pIdI, err := strconv.Atoi(pId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Product ID"})
			c.Abort()
			return
		}
		product, err := store.GetProductById(c.Request.Context(), pIdI)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Product with id of %d not found", pIdI)})
				c.Abort()
				return
			}
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			c.Abort()
			return
		}
		if product.UserId != user.Id {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "You are forbidden from performing this operation"})
			c.Abort()
			return
		}
		c.Set("product", product)
		c.Set("id", pIdI)
		c.Next()
	}
}
