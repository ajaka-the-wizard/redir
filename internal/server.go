package internal

import (
	"context"

	"github.com/ajaka-the-wizard/redir/internal/cache"
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/database"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/ajaka-the-wizard/redir/internal/routes"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/gin-gonic/gin"
)

func Listen() error {
	ctx := context.Background()
	cfg := configs.LoadEnv()
	rdb := cache.InitializeRedis(ctx, cfg)
	defer rdb.Clean()
	store := store.InitializeStore(rdb)
	_, presignedClient, tm := configs.PerformAllNecessaryActivationStep(ctx, cfg)

	pool := database.ConnectDB(ctx, cfg.DATABASE_URL)
	defer pool.Close()

	if cfg.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middlewares.GenAndAttachRequestIdMiddleware())
	router.Use(middlewares.AttachLoggerToContext())
	router.Use(middlewares.PerformBasicRequestCycleCalculations())
	router.Use(gin.Recovery())

	router.SetTrustedProxies(nil)

	v1 := router.Group("/api/v1")

	routes.AuthRoutes(v1, pool, cfg, store)
	routes.UserRoutes(v1, pool, cfg, store)
	routes.ProductRoutes(v1, pool, cfg, store)
	routes.ClientRoutes(v1, pool, cfg, tm, store)
	routes.AssetRoutes(v1, pool, cfg, presignedClient, store)

	return router.Run(cfg.SERVER_ADDRESS)
}
