package handlers

import (
	"log"
	"net/http"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/gin-gonic/gin"
)

func GenerateKey(cfg *configs.EnvData, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user, ok := utils.GetUser(c)
		pIdI, ok := utils.GetId(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Couldn't get id"})
			return
		}
		pKey := utils.GeneratePrivateKey()
		hKey, err := utils.PerformMultiStepHash(pKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		product, err := store.CreatePrivateKey(c.Request.Context(), pIdI, hKey)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		product.PrivateKey = pKey
		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Private Key created successfully", "product": product})
	}
}

func CreateProduct(cfg *configs.EnvData, store *store.Store) gin.HandlerFunc {
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
		product, err := store.CreateProduct(c.Request.Context(), request)
		if err != nil {
			log.Println("CP")
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Product created successfully", "product": &product})
	}
}
