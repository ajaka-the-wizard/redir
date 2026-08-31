package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func InitializeRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool,
	}
}

func (r *Repository) CheckHealth(ctx context.Context, logger *slog.Logger) error {
	logger.Info("postgres: checking readiness")

	if r == nil || r.pool == nil {
		return fmt.Errorf("postgres: pool is nil")
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.pool.Ping(ctx); err != nil {
		logger.Error("postgres: readiness check failed", "error", err.Error())
		return fmt.Errorf("postgres: %w", err)
	}

	logger.Info("postgres: readiness check passed")
	return nil
}
