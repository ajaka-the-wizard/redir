package internal

import (
	"context"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/database"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/ajaka-the-wizard/redir/internal/routes"
	"github.com/gin-gonic/gin"
)

func Listen() error {
	ctx := context.Background()
	cfg := configs.LoadEnv()

	client, presignedClient := configs.PerformAllNecessaryActivationStep(ctx, cfg)

	pool := database.ConnectDB(ctx, cfg.DATABASE_URL)
	mmap := memory.NewMemoryMap()

	if cfg.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}
	defer pool.Close()

	router := gin.New()

	router.Use(middlewares.GenAndAttachRequestIdMiddleware())
	router.Use(middlewares.AttachLoggerToContext())
	router.Use(middlewares.PerformBasicRequestCycleCalculations())
	router.Use(gin.Recovery())

	router.SetTrustedProxies(nil)

	v1 := router.Group("/api/v1")

	routes.AuthRoutes(v1, pool, cfg, mmap)
	routes.UserRoutes(v1, pool, cfg, mmap)
	routes.ProductRoutes(v1, pool, cfg, mmap)
	routes.ClientRoutes(v1, pool, cfg, client)
	routes.AssetRoutes(v1, pool, cfg, presignedClient)

	return router.Run(cfg.SERVER_ADDRESS)
}
