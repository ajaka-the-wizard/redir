package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
)

func HandleRedirect(cfg *configs.EnvData, presignedClient *s3.PresignClient, store *store.Store, parser *uaparser.Parser) gin.HandlerFunc {
	return func(c *gin.Context) {
		media, ok := utils.GetMedia(c)
		ua := c.GetHeader("user-agent")
		log.Println("parsing", ua)
		client := parser.Parse(ua)
		metric := models.Metrics{
			MediaId:        media.PublicKey,
			Browser:        client.UserAgent.Family,
			BrowserVersion: client.UserAgent.Major,
			Device:         client.Device.Family,
			DeviceBrand:    client.Device.Brand,
			DeviceModel:    client.Device.Model,
			Os:             client.Os.Family,
			OsVersion:      client.Os.Major,
			Ip:             c.ClientIP(),
		}
		log.Println("metrics:", metric)
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
		err = store.SaveMetrics(c.Request.Context(), &metric)
		if err != nil {
			log.Println("Error saving metric:", err)
		}
		c.Redirect(http.StatusFound, preSignedUrl)
	}
}
