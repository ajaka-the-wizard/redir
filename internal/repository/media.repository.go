package repository

import (
	"context"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateMedia(ctx context.Context, pool *pgxpool.Pool, user_id uuid.UUID, innerKey string, publicKey string) (*models.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
	INSERT INTO medias (public_key, inner_key,user_id)
	VALUES ($1, $2, $3)
	RETURNING public_key, user_id, file_size, status, file_type, active, file_name, created_at, updated_at
	`
	rows, err := pool.Query(ctx, query, publicKey, innerKey)
	if err != nil {
		return nil, err
	}
	media, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Media])
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func CreateMediaBatch(pool *pgxpool.Pool, mediaBatch *[]models.Media) *[]models.Media {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	var allSavedMedia []models.Media
	batch := pgx.Batch{}

	for _, f := range *mediaBatch {
		batch.Queue(
			`
				INSERT INTO medias (public_key, inner_key, user_id, file_name, mime_type, batch_id, seq_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING public_key, user_id, file_size, status, file_type, active, file_name, created_at, updated_at
				`, f.PublicKey, f.InnerKey, f.UserId, f.FileName, f.MimeType, f.BatchId, f.SeqId,
		)
	}
	results := pool.SendBatch(ctx, &batch)
	defer results.Close()

	for i := 0; i < len(*mediaBatch); i++ {
		rows, err := results.Query()

		if err != nil {
			continue
		}
		m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Media])
		if err != nil {
			continue
		}
		allSavedMedia = append(allSavedMedia, m)
	}
	return &allSavedMedia
}

func HandleBatchCommits(ctx context.Context, pool *pgxpool.Pool, batchId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	UPDATE medias
	SET status = completed
	WHERE batch_id = $1
	RETURNING public_key
	`
	err := pool.QueryRow(ctx, query, batchId).Scan()
	if err != nil {
		return err
	}
	return nil
}

func RetriveBatch(ctx context.Context, pool *pgxpool.Pool, batchId uuid.UUID) (*[]models.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	SELECT * 
	FROM medias
	WHERE batch_id = $1
	`
	rows, err := pool.Query(ctx, query, batchId)
	if err != nil {
		return nil, err
	}
	medias, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Media])
	if err != nil {
		return nil, err
	}
	return &medias, nil
}

func GetMedia(ctx context.Context, pool *pgxpool.Pool, publicKey string) (*models.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	SELECT *
	FROM medias
	WHERE public_key = $1
	`
	rows, err := pool.Query(ctx, query, publicKey)
	if err != nil {
		return nil, err
	}
	media, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Media])
	if err != nil {
		return nil, err
	}
	return &media, nil
}
