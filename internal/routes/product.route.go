package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, cfg *configs.EnvData, store *store.Store) {
	product := rg.Group("/product")
	product.Use(middlewares.RL.GetLimiterForProductAndUser(10))
	product.Use(middlewares.AuthMiddleware(store, cfg))
	product.POST("", middlewares.ProductValidationMiddleware, handlers.CreateProduct(pool, cfg, store))
	product.POST("/:id", middlewares.CanThisUserAlterThisProduct(pool, cfg, store), handlers.GenerateKey(pool, cfg, store))
}
