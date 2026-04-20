package handlers

import (
	"net/http"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func HandleRedirect(cfg *configs.EnvData, presignedClient *s3.PresignClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		media, ok := utils.GetMedia(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}

		preSignedUrl, err := presignedClient.PresignGetObject(c.Request.Context(), &s3.GetObjectInput{
			Bucket: aws.String(cfg.BUCKET_NAME),
			Key:    aws.String(media.InnerKey),
		}, s3.WithPresignExpires(time.Minute*30))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		c.Redirect(http.StatusFound, preSignedUrl.URL)
	}
}
