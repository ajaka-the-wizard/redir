package store

import (
	"context"
	"log/slog"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/domain"
)

func (s *Store) SetUserOnline(ctx context.Context, logger *slog.Logger, sessionId string, u *domain.
	LightUser) (time.Time, error) {
	return s.r.SetUserOnline(ctx, sessionId, u)
}

func (s *Store) GetUser(ctx context.Context, logger *slog.Logger, sessionId string) (*domain.
	LightUser, bool) {
	return s.r.GetUser(ctx, sessionId)
}

func (s *Store) RevokeUser(ctx context.Context, logger *slog.Logger, sessionId string) error {
	return s.r.RevokeUser(ctx, sessionId)
}
