package internal

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/database"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/ajaka-the-wizard/redir/internal/routes"
	"github.com/gin-gonic/gin"
)

func Listen() {
	cfg := configs.LoadEnv()
	pool := database.Connect_DB(cfg.DATABASEURL)
	mmap := memory.NewMemoryMap()
	defer pool.Close()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	v1 := router.Group("/api/v1")
	routes.AuthRoutes(v1, pool, cfg, mmap)
	routes.UserRoutes(v1, pool, cfg, mmap)
	routes.ProductRoutes(v1, pool, cfg, mmap)
	routes.ClientRoutes(v1, pool, cfg)
	router.Run(cfg.SERVERADDRESS)
}
