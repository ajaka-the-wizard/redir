package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/gin-gonic/gin"
)

func ClientRoutes(rg *gin.RouterGroup, cfg *configs.EnvData, tm *transfermanager.Client, store *store.Store) {
	clients := rg.Group("/client")
	clients.Use(middlewares.RL.GetLimiterForClient(15))
	clients.Use(middlewares.CheckAndValidateClientKeys(cfg, store))
	clients.GET("/ping", handlers.HandleClientPing())
	clients.POST("/upload", handlers.HandleUpload(cfg, tm, store))
	clients.PUT("/commit/:batchId", handlers.HandleBatchCommit(store))
}
