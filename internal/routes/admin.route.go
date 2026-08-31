package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(rg *gin.RouterGroup, s *store.Store) {
	rg.GET("/readyz", handlers.HandleReady(s))
}
