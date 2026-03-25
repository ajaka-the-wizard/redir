package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ClientRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, cfg *configs.EnvData) {
	client := rg.Group("/client")
	client.Use(middlewares.RL.GetLimiterForClient(15))
	client.Use(middlewares.CheckAndValidateClientKeys(pool, cfg))
	client.GET("/ping", handlers.HandleClientPing())
	client.POST("/upload", handlers.HandleUpload())
}
