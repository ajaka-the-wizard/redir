package store

import (
	"context"
	"log"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (s *Store) CreateMedia(ctx context.Context, pool *pgxpool.Pool, user_id uuid.UUID, innerKey string, publicKey string) (*models.Media, error) {
	return repository.CreateMedia(ctx, pool, user_id, innerKey, publicKey)
}

func (s *Store) CreateMediaBatch(ctx context.Context, pool *pgxpool.Pool, mediaBatch *[]models.Media) *[]models.Media {
	return repository.CreateMediaBatch(ctx, pool, mediaBatch)
}

func (s *Store) HandleBatchCommits(ctx context.Context, pool *pgxpool.Pool, batchId uuid.UUID) error {
	return repository.HandleBatchCommits(ctx, pool, batchId)
}

func (s *Store) RetriveBatch(ctx context.Context, pool *pgxpool.Pool, batchId uuid.UUID) (*[]models.Media, error) {
	return repository.RetriveBatch(ctx, pool, batchId)
}

func (s *Store) GetMedia(ctx context.Context, pool *pgxpool.Pool, publicKey string) (*models.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	m, err := s.r.GetMedia(ctx, publicKey)
	if err == nil {
		return m, nil
	}
	m, err = repository.GetMedia(ctx, pool, publicKey)
	if err != nil {
		return nil, err
	}
	err = s.r.SetMedia(ctx, *m)
	if err != nil {
		log.Println("couldnt set", err)
	}
	return m, nil
}

func (s *Store) SetPresigned(ctx context.Context, publicKey string, url string, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.r.SetPresignedUrl(ctx, publicKey, url, exp)
}

func (s *Store) GetPresigned(ctx context.Context, publicKey string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.r.GetPresignedUrl(ctx, publicKey)
}
