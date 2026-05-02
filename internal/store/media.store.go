package store

import (
	"context"
	"log"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/google/uuid"
)

func (s *Store) CreateMedia(ctx context.Context, user_id uuid.UUID, innerKey string, publicKey string) (*models.Media, error) {
	return s.repo.CreateMedia(ctx, user_id, innerKey, publicKey)
}

func (s *Store) CreateMediaBatch(ctx context.Context, mediaBatch *[]models.Media) *[]models.Media {
	return s.repo.CreateMediaBatch(ctx, mediaBatch)
}

func (s *Store) HandleBatchCommits(ctx context.Context, batchId uuid.UUID) error {
	return s.repo.HandleBatchCommits(ctx, batchId)
}

func (s *Store) RetriveBatch(ctx context.Context, batchId uuid.UUID) (*[]models.Media, error) {
	return s.repo.RetriveBatch(ctx, batchId)
}

func (s *Store) GetMedia(ctx context.Context, publicKey string) (*models.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	m, err := s.r.GetMedia(ctx, publicKey)
	if err == nil {
		return m, nil
	}
	m, err = s.repo.GetMedia(ctx, publicKey)
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
