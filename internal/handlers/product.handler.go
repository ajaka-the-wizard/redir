package handlers

import (
	"log"
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GenerateKey(pool *pgxpool.Pool, cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user, ok := utils.GetUser(c)
		pIdI, ok := utils.GetId(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Couldn't get id"})
			return
		}
		p_key := utils.GeneratePrivateKey()
		h_key, err := utils.PerformMultiStepHash(p_key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		product, err := repository.CreatePrivateKey(pool, cfg, pIdI, h_key)
		if err != nil {
			log.Println("Hey")
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		product.PrivateKey = p_key
		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Private Key created successfully", "product": product})
	}
}

func CreateProduct(pool *pgxpool.Pool, cfg *configs.EnvData) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := utils.GetUser(c)
		val, _ := c.Get("product")
		request, ok := val.(*domain.CreateProductDetails)
		if !ok {
			log.Println("product")
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		request.UserId = user.Id
		product, err := repository.CreateProduct(pool, cfg, request)
		if err != nil {
			log.Println("CP")
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Product created successfully", "product": &product})
	}
}
