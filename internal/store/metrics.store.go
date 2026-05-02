package store

import (
	"context"

	"github.com/ajaka-the-wizard/redir/internal/models"
)

func (s *Store) SaveMetrics(ctx context.Context, metric *models.Metrics) error {
	return s.repo.SaveMetrics(ctx, metric)
}
