package internal

import (
	"context"
	"log/slog"

	"github.com/ajaka-the-wizard/redir/internal/cache"
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/database"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/routes"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/gin-gonic/gin"
)

func Listen() error {
	logger := slog.Default()
	ctx := context.Background()
	cfg := configs.LoadEnv(logger)
	parser := configs.InitializeUserAgentParser()
	rdb := cache.InitializeRedis(ctx, cfg, logger)
	defer rdb.Clean()
	_, presignedClient, tm := configs.PerformAllNecessaryActivationStep(ctx, cfg, logger)
	pool := database.ConnectDB(ctx, logger, cfg.DATABASE_URL)
	defer pool.Close()
	repo := repository.InitializeRepository(pool)
	store := store.InitializeStore(rdb, repo)
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
	routes.AuthRoutes(v1, cfg, store)
	routes.UserRoutes(v1, cfg, store)
	routes.ProductRoutes(v1, cfg, store)
	routes.ClientRoutes(v1, cfg, tm, store)
	routes.AssetRoutes(v1, cfg, presignedClient, store, parser)
	return router.Run(cfg.SERVER_ADDRESS)
}
