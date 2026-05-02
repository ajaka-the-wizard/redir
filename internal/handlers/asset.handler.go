package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func HandleRedirect(cfg *configs.EnvData, presignedClient *s3.PresignClient, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		media, ok := utils.GetMedia(c)
		d := c.GetHeader("user-agent")
		log.Println("d", d)
		f := c.GetHeader("User-Agent")
		log.Println("f", f)
		if !ok {
			log.Println("From ok")
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		preSignedUrl, err := store.GetPresigned(c.Request.Context(), media.PublicKey)
		if err != nil {
			preSigned, err := presignedClient.PresignGetObject(c.Request.Context(), &s3.GetObjectInput{
				Bucket: aws.String(cfg.BUCKET_NAME),
				Key:    aws.String(media.InnerKey),
			}, s3.WithPresignExpires(30*time.Minute))
			if err != nil {
				log.Println("from presigned", err)
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
				return
			}
			store.SetPresigned(c.Request.Context(), media.PublicKey, preSigned.URL, 28*time.Minute)
			preSignedUrl = preSigned.URL
		}
		c.Redirect(http.StatusFound, preSignedUrl)
	}
}
