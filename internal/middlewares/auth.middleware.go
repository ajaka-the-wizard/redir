package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		c.Set("sessionId", sessionId)
		c.Set("user", user)
		c.Next()
	}
}

func CheckAndValidateClientKeys(pool *pgxpool.Pool, cfg *configs.EnvData) gin.HandlerFunc {
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
		key, err := repository.GetProductById(pool, cfg, pIdI)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Product with id of %d not found", pIdI)})
				c.Abort()
				return
			}
			log.Println("xxx")
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			c.Abort()
			return
		}
		err = utils.VerifyMultipStepHash(k, key.PrivateKey)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid key given"})
			c.Abort()
			return
		}
		c.Set("client", key)
		c.Next()
	}
}

func CanThisUserAlterThisProduct(pool *pgxpool.Pool, cfg *configs.EnvData) gin.HandlerFunc {
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
		product, err := repository.GetProductById(pool, cfg, pIdI)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Product with id of %d not found", pIdI)})
				c.Abort()
				return
			}
			log.Println("xxx")
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
		c.Set("product", &product)
		c.Set("id", pIdI)
		c.Next()
	}
}
