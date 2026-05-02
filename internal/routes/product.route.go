package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(rg *gin.RouterGroup, cfg *configs.EnvData, store *store.Store) {
	product := rg.Group("/product")
	product.Use(middlewares.RL.GetLimiterForProductAndUser(10))
	product.Use(middlewares.AuthMiddleware(store, cfg))
	product.POST("", middlewares.ProductValidationMiddleware, handlers.CreateProduct(cfg, store))
	product.POST("/:id", middlewares.CanThisUserAlterThisProduct(cfg, store), handlers.GenerateKey(cfg, store))
}
