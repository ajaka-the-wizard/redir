package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, cfg *configs.EnvData, mmap *memory.AuthMemoryMap) {
	product := rg.Group("/product")
	product.Use(middlewares.RL.GetLimiterForProductAndUser(10))
	product.Use(middlewares.AuthMiddleware(mmap))
	product.POST("", middlewares.ProductValidationMiddleware, handlers.CreateProduct(pool, cfg))
	product.POST("/:id", middlewares.CanThisUserAlterThisProduct(pool, cfg), handlers.GenerateKey(pool, cfg))
}
