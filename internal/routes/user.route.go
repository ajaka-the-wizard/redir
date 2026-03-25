package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, cfg *configs.EnvData, mmap *memory.AuthMemoryMap) {
	user := rg.Group("/users")
	user.Use(middlewares.RL.GetLimiterForProductAndUser(10))
	user.Use(middlewares.AuthMiddleware(mmap))
	user.GET("/me", handlers.GetUser(pool, cfg))
}
