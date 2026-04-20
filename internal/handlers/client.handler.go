package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleClientPing() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := utils.GetLogger(c)
		logger.Info("ping attempt from client")
		client, ok := utils.GetProduct(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Client %v recognised, Proxy is ready for you.", client.ProductId)})
		logger.Info("Client ping ponged", "productId", client.ProductId)
	}
}

func HandleUpload(cfg *configs.EnvData, pool *pgxpool.Pool, client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var mediaBatch []models.Media
		product, ok := utils.GetProduct(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		batchId := c.PostForm("batch_id")
		batchIdUUID, err := utils.ValidateAndReturnUUID(batchId)
		if batchId == "" || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "'batch_id' is either missing or not valid uuid"})
			return
		}
		reader, err := c.Request.MultipartReader()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Error getting multipart reader"})
			return
		}
		medias, err := repository.RetriveBatch(c.Request.Context(), pool, batchIdUUID)
		if err != nil && err != pgx.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		existing := map[int]struct{}{}
		if medias != nil {
			for _, m := range *medias {
				existing[m.SeqId] = struct{}{}
			}
		}

		for {
			part, err := reader.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				}
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "multipart read error"})
				return

			}
			seqId := part.Header.Get("X-Sequential-ID")
			seqIdI, _ := strconv.Atoi(seqId)

			if _, ok := existing[seqIdI]; ok {
				continue
			}

			fileName := part.FileName()
			contentType := part.Header.Get("Content-Type")
			innerKey, publicKey := utils.GenerateKeyForUpload(cfg, product.ProductId)
			if _, err := client.PutObject(c.Request.Context(), &s3.PutObjectInput{
				Bucket: &cfg.BUCKET_NAME,
				Body:   part,
				Key:    &innerKey,
				Metadata: map[string]string{
					"original_name": fileName,
					"content_type":  contentType,
				},
			}); err != nil {
				continue
			}
			mediaBatch = append(mediaBatch, models.Media{
				PublicKey: publicKey,
				InnerKey:  innerKey,
				FileName:  fileName,
				MimeType:  contentType,
				UserId:    product.UserId,
				BatchId:   batchIdUUID,
				SeqId:     seqIdI,
				Public:    product.Public,
			})
		}
		media := repository.CreateMediaBatch(c.Request.Context(), pool, &mediaBatch)
		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "file uploaded successfully", "media": media})
	}
}

func HandleBatchCommit(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		batchId := c.Param("batchId")
		batchIdUUID, err := utils.ValidateAndReturnUUID(batchId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid id given"})
			return
		}
		err = repository.HandleBatchCommits(c.Request.Context(), pool, batchIdUUID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": fmt.Sprintf("Could not find batch with Id of '%s'", batchIdUUID)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Batch commited successfully"})
	}
}
