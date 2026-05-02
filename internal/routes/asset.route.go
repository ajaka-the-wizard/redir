package routes

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/handlers"
	"github.com/ajaka-the-wizard/redir/internal/middlewares"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
)

func AssetRoutes(rg *gin.RouterGroup, cfg *configs.EnvData, presignedClient *s3.PresignClient, store *store.Store, parser *uaparser.Parser) {
	asset := rg.Group("/assets")

	asset.GET("/:assetId", middlewares.CheckIfAssetIsPublic(cfg, store), handlers.HandleRedirect(cfg, presignedClient, store, parser))
}
