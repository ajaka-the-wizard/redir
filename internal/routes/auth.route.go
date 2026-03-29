package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/memory"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, cfg *configs.EnvData, mmap *memory.AuthMemoryMap) {
	auth := rg.Group("/auth")
	auth.POST("/register", middlewares.RegisterValidationMiddleware, handlers.HandleRegister(pool, cfg))
	auth.POST("/login", middlewares.LoginValidationMiddleware, handlers.HandleLogin(pool, cfg, mmap))
	auth.POST("/logout", middlewares.AuthMiddleware(mmap), handlers.HandleLogout(mmap, cfg))
	o := handlers.InitGoogleOauth(cfg)
	og := handlers.InitGithubOauth(cfg)
	auth.GET("/oauth/google", o.HandleRedirectToGoogle(cfg))
	auth.GET("/oauth/github", og.HandleRedirectToGithub(cfg))
	auth.GET("oauth/google/callback", o.HandleGoogleCallback(pool, cfg, mmap))
}
