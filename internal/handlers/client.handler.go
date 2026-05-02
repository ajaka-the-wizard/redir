package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/ajaka-the-wizard/redir/internal/store"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func HandleUpload(cfg *configs.EnvData, tm *transfermanager.Client, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var mediaBatch []models.Media
		valid := false
		product, ok := utils.GetProduct(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
			return
		}
		batchId := c.GetHeader("X-Batch-ID")
		batchIdUUID, err := utils.ValidateAndReturnUUID(batchId)
		if batchId == "" || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "'X-Batch-ID' header is either missing or not valid uuid"})
			return
		}
		reader, err := c.Request.MultipartReader()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Error getting multipart reader"})
			return
		}
		medias, err := store.RetriveBatch(c.Request.Context(), batchIdUUID)
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
					if !valid {
						c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Nothing was received"})
						return
					}
					break
				}
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "multipart read error"})
				return
			}
			valid = true
			seqId := part.Header.Get("X-Sequential-ID")
			seqIdI, _ := strconv.Atoi(seqId)
			log.Println("Checking existence", part.FileName(), seqIdI)
			if _, ok := existing[seqIdI]; ok {
				log.Println("found", part.FileName(), seqIdI)
				continue
			}

			fileName := part.FileName()
			contentType := part.Header.Get("Content-Type")
			innerKey, publicKey := utils.GenerateKeyForUpload(cfg, product.ProductId)
			log.Println("Uploading file", part.FileName())
			_, err = tm.UploadObject(c.Request.Context(), &transfermanager.UploadObjectInput{
				Bucket: &cfg.BUCKET_NAME,
				Key:    &innerKey,
				Body:   part,
				Metadata: map[string]string{
					"original_name": fileName,
					"content_type":  contentType,
				},
			})

			if err != nil {
				log.Println("Error while uploading", seqIdI, err)
				continue

			}
			log.Println("Batching", seqIdI)
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
		if len(mediaBatch) == 0 {
			c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "They all failed"})
			return
		}
		log.Println("batch", len(mediaBatch))
		media := store.CreateMediaBatch(c.Request.Context(), &mediaBatch)
		HydrateMedias(cfg, *media)
		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "file uploaded successfully", "media": media})
	}
}

func HandleBatchCommit(store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		batchId := c.Param("batchId")
		batchIdUUID, err := utils.ValidateAndReturnUUID(batchId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid id given"})
			return
		}
		err = store.HandleBatchCommits(c.Request.Context(), batchIdUUID)

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
