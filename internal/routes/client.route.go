package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ClientRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, cfg *configs.EnvData, tm *transfermanager.Client) {
	clients := rg.Group("/client")
	clients.Use(middlewares.RL.GetLimiterForClient(15))
	clients.Use(middlewares.CheckAndValidateClientKeys(pool, cfg))
	clients.GET("/ping", handlers.HandleClientPing())
	clients.POST("/upload", handlers.HandleUpload(cfg, pool, tm))
	clients.PUT("/commit/:batchId", handlers.HandleBatchCommit(pool))
}
