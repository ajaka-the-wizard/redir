package repository

import (
	"context"
	"log"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateMedia(ctx context.Context, user_id uuid.UUID, innerKey string, publicKey string) (*models.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
	INSERT INTO medias (public_key, inner_key,user_id)
	VALUES ($1, $2, $3)
	RETURNING public_key, inner_key, user_id, file_size, status, file_type, active, file_name, created_at, updated_at, batch_id,seq_id,public,mime_type,hits,
	`
	rows, err := r.pool.Query(ctx, query, publicKey, innerKey, user_id)
	if err != nil {
		return nil, err
	}
	media, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Media])
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *Repository) CreateMediaBatch(ctx context.Context, mediaBatch *[]models.Media) *[]models.Media {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var allSavedMedia []models.Media
	batch := pgx.Batch{}

	for _, f := range *mediaBatch {
		batch.Queue(
			`
				INSERT INTO medias (public_key, inner_key, user_id, file_name, mime_type, batch_id, seq_id, public)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING public_key, inner_key, user_id, file_size, status, file_type, active, file_name, created_at, updated_at, batch_id,seq_id, public, mime_type, hits
				`, f.PublicKey, f.InnerKey, f.UserId, f.FileName, f.MimeType, f.BatchId, f.SeqId, f.Public,
		)
	}
	results := r.pool.SendBatch(ctx, &batch)
	defer results.Close()

	for range *mediaBatch {
		rows, err := results.Query()

		if err != nil {
			log.Println("error with query", err)
			continue
		}
		m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Media])
		if err != nil {
			log.Println("error with collect", err)
			continue
		}
		allSavedMedia = append(allSavedMedia, m)
	}
	return &allSavedMedia
}

func (r *Repository) HandleBatchCommits(ctx context.Context, batchId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	UPDATE medias
	SET status = 'completed'
	WHERE batch_id = $1
	`
	tag, err := r.pool.Exec(ctx, query, batchId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *Repository) RetriveBatch(ctx context.Context, batchId uuid.UUID) (*[]models.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	SELECT * 
	FROM medias
	WHERE batch_id = $1 and status = 'pending'
	`
	rows, err := r.pool.Query(ctx, query, batchId)
	if err != nil {
		return nil, err
	}
	medias, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Media])
	if err != nil {
		return nil, err
	}
	return &medias, nil
}

func (r *Repository) GetMedia(ctx context.Context, publicKey string) (*models.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	SELECT *
	FROM medias
	WHERE public_key = $1
	`
	rows, err := r.pool.Query(ctx, query, publicKey)
	if err != nil {
		return nil, err
	}
	media, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Media])
	if err != nil {
		return nil, err
	}
	return &media, nil
}
