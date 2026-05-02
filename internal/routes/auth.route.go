package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, cfg *configs.EnvData, store *store.Store) {
	auth := rg.Group("/auth")

	o := handlers.InitGoogleOauth(cfg)
	og := handlers.InitGithubOauth(cfg)

	auth.Use(middlewares.RL.GetLimiterForAuth(5))

	auth.POST("/register", middlewares.RegisterValidationMiddleware, handlers.HandleRegister(cfg, store))
	auth.POST("/login", middlewares.LoginValidationMiddleware, handlers.HandleLogin(cfg, store))
	auth.POST("/logout", middlewares.AuthMiddleware(store, cfg), handlers.HandleLogout(store, cfg))
	auth.GET("/oauth/google", o.HandleRedirectToGoogle(cfg))
	auth.GET("/oauth/github", og.HandleRedirectToGithub(cfg))
	auth.GET("/oauth/google/callback", o.HandleGoogleCallback(cfg, store))
}
